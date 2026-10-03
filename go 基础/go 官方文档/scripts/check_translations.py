#!/usr/bin/env python3
"""Check translated documents against the preserved English source snapshot."""

import argparse
from collections import Counter
import hashlib
import html
import json
from pathlib import Path
import re
import subprocess
import textwrap

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / "translations/en"
CONTENT = ROOT / "_content"
MANIFEST = ROOT / "translations/manifest.json"


def code_blocks(text):
    blocks = re.findall(r"<pre\b[^>]*>(.*?)</pre\s*>", text, re.S | re.I)
    # Copy buttons live inside some <pre> elements; their translated UI labels
    # are not part of the executable example.
    blocks = [html.unescape(re.sub(r"<button\b.*?</button\s*>", "", block, flags=re.S | re.I)) for block in blocks]
    fence = None
    current = []
    for line in text.splitlines():
        match = re.match(r"^\s*(`{3,}|~{3,})(.*)$", line)
        if fence is None:
            if match:
                fence = match[1]
                current = []
        elif match and match[1][0] == fence[0] and len(match[1]) >= len(fence) and not match[2].strip():
            blocks.append("\n".join(current))
            fence = None
        else:
            current.append(line)
    if fence:
        # Markdown permits an EOF-terminated fenced block (used by upstream).
        blocks.append("\n".join(current))
    return [textwrap.dedent(block).strip("\n") for block in blocks]


def anchors(text):
    return set(re.findall(r'\bid=["\']([^"\']+)["\']|\{#([^}]+)\}', text))


def api_headings(text):
    headings = {}
    for match in re.finditer(r'<h([23])\b[^>]*\bid="([^"]+)"[^>]*>(.*?)</h\1>', text, re.S):
        if re.match(r"(?:func|type)\s", match[3]):
            headings[match[2]] = re.findall(r"<a\b[^>]*>(.*?)</a>", match[3], re.S)
    return headings


def api_example_names(text):
    # Example suffixes describe a scenario and may be translated; the API name
    # before them must remain intact. Names starting with '_' are package examples.
    return {name: re.split(r"[（(]", label)[0].strip()
            for name, label in re.findall(r'<a class="exampleLink" href="#example_([^"]*)">([^<]*)</a>', text)
            if name and not name.startswith("_")}


def rendered_snippets(text, tag):
    blocks = re.findall(rf"<{tag}\b[^>]*>(.*?)</{tag}\s*>", text, re.S | re.I)
    if tag == "pre":
        blocks = [re.sub(r"<[^>]*>", "", block) for block in blocks]
    return [html.unescape(block) for block in blocks]


def check_code(path, original, translated, prose_blocks, comment_blocks=None, comment_checks=None):
    before, after = code_blocks(original), code_blocks(translated)
    if len(before) != len(after):
        return [f"{path}: Number of code blocks changed"]
    exceptions = {block["index"]: block for block in prose_blocks.get(path, [])}
    comments = {(block["kind"], block["index"]): block for block in (comment_blocks or {}).get(path, [])}
    used_comments = set()
    failures = []

    def check_comment(kind, index, source, target):
        record = comments.get((kind, index))
        if not record or comment_checks is None:
            return False
        used_comments.add((kind, index))
        if (hashlib.sha256(source.encode()).hexdigest() != record["source_sha256"]
                or hashlib.sha256(target.encode()).hexdigest() != record["translation_sha256"]):
            failures.append(f"{path}: Reviewed {kind} comments {index} changed")
        before_plain = rendered_snippets(original, kind)[index]
        after_plain = rendered_snippets(translated, kind)[index]
        key = f"{path}:{kind}:{index}"
        comment_checks[key + ":source"] = before_plain
        comment_checks[key + ":translation"] = after_plain
        return True

    for i, (source, target) in enumerate(zip(before, after)):
        if i in exceptions:
            block = exceptions[i]
            if (hashlib.sha256(source.encode()).hexdigest() != block["source_sha256"]
                    or hashlib.sha256(target.encode()).hexdigest() != block["translation_sha256"]):
                failures.append(f"{path}: Reviewed prose table {i} changed")
        elif source != target and not check_comment("pre", i, source, target):
            failures.append(f"{path}: Code block {i} differs from source")
    if exceptions.keys() - set(range(len(before))):
        failures.append(f"{path}: Prose table index out of range")
    pattern = r"<textarea\b[^>]*>(.*?)</textarea\s*>"
    before_run, after_run = re.findall(pattern, original, re.S), re.findall(pattern, translated, re.S)
    if len(before_run) != len(after_run):
        failures.append(f"{path}: Number of runnable examples changed")
    for i, (source, target) in enumerate(zip(before_run, after_run)):
        if source != target and not check_comment("textarea", i, source, target):
            failures.append(f"{path}: Runnable example {i} differs from source")
    if comments.keys() - used_comments:
        failures.append(f"{path}: Stale or out-of-range reviewed comment records")
    return failures


def check_partial(path, original, translated, selected, occurrences=None):
    pattern = r'<h([2-6])\b[^>]*\bid="([^"]+)"[^>]*>.*?</h\1>'
    before = list(re.finditer(pattern, original, re.S))
    after = list(re.finditer(pattern, translated, re.S))
    if [m[2] for m in before] != [m[2] for m in after]:
        return [f"{path}: Partial translation heading order changed"]
    ranges = []
    seen = Counter()
    for i, heading in enumerate(before):
        seen[heading[2]] += 1
        allowed = (occurrences or {}).get(heading[2])
        if heading[2] in selected and (allowed is None or seen[heading[2]] in allowed):
            end = next((h.start() for h in before[i + 1:] if int(h[1]) <= int(heading[1])), len(original))
            ranges.append((heading.start(), end))
    failures = []
    for i, heading in enumerate(before):
        if any(a <= heading.start() < b for a, b in ranges):
            continue
        source_end = before[i + 1].start() if i + 1 < len(before) else len(original)
        target_end = after[i + 1].start() if i + 1 < len(after) else len(translated)
        if original[heading.start():source_end] != translated[after[i].start():target_end]:
            failures.append(f"{path}: Unselected section {heading[2]} changed")
    return failures


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base-url", help="Also check HTTP rendering on a running local website")
    args = parser.parse_args()
    manifest = json.loads(MANIFEST.read_text(encoding="utf-8"))
    prose_blocks = json.loads((ROOT / "translations/prose-blocks.json").read_text())
    comment_blocks = json.loads((ROOT / "translations/comment-blocks.json").read_text())
    comment_checks = {}
    failures = []
    count = 0
    markdown_paths = []
    for entry in manifest["documents"]:
        if entry["status"] == "excluded-history":
            path = entry["path"]
            if (SOURCE / path).read_bytes() != (CONTENT / path).read_bytes():
                failures.append(f"{path}: Excluded release notes changed")
        if entry["status"] not in {"translated", "partial"}:
            continue
        path = entry["path"]
        original = (SOURCE / path).read_text(encoding="utf-8")
        translated = (CONTENT / path).read_text(encoding="utf-8")
        if path.endswith(".md"):
            markdown_paths.extend([str(SOURCE / path), str(CONTENT / path)])
        if hashlib.sha256(original.encode()).hexdigest() != entry["source_sha256"]:
            failures.append(f"{path}: English source hash changed")
        if not re.search(r"[\u4e00-\u9fff]", translated):
            failures.append(f"{path}: Chinese text missing")
        failures.extend(check_code(path, original, translated, prose_blocks, comment_blocks, comment_checks))
        if anchors(original) - anchors(translated):
            failures.append(f"{path}: Original section anchors missing")
        if path.startswith("pkg/") and api_headings(original) != api_headings(translated):
            failures.append(f"{path}: API names in declaration headings changed")
        if path.startswith("pkg/") and api_example_names(original) != api_example_names(translated):
            failures.append(f"{path}: API names in example navigation changed")
        if path.startswith(("pkg/", "cmd/", "ref/")) and path.endswith(".html"):
            links = lambda text: Counter(re.findall(r'\b(?:href|src)=["\']([^"\']+)["\']', text))
            if links(original) - links(translated):
                failures.append(f"{path}: Original reference links missing or changed")
        if entry["status"] == "partial":
            failures.extend(check_partial(path, original, translated, entry["coverage"]["sections"],
                                          entry["coverage"].get("section_occurrences")))
        if path.endswith(".article"):
            pattern = r"(?m)^\.(?:play|code|image|iframe|html|background)\b.*$"
            if re.findall(pattern, original) != re.findall(pattern, translated):
                failures.append(f"{path}: Tour directives changed")
        if args.base_url:
            import urllib.request
            url = args.base_url.rstrip("/") + entry["url"]
            if path.endswith(".article"):
                url = args.base_url.split("/go.dev")[0].rstrip("/") + "/tour/lesson/" + Path(path).stem
            try:
                with urllib.request.urlopen(url, timeout=15) as response:
                    body = response.read().decode("utf-8")
                    if response.status != 200 or "[template error:" in body or "加载文件失败：" in body:
                        failures.append(f"{path}: Render failed ({response.status})")
                    if not re.search(r"[\u4e00-\u9fff]", body):
                        failures.append(f"{path}: Served content is not Chinese")
            except Exception as error:
                failures.append(f"{path}: HTTP check failed: {error}")
        count += 1
    if comment_checks:
        try:
            result = subprocess.run(
                ["go", "run", str(ROOT / "scripts/go_comment_structure.go")],
                input=json.dumps(comment_checks), cwd=ROOT, check=True, capture_output=True, text=True)
            structures = json.loads(result.stdout)
            for key in comment_checks:
                if key.endswith(":source"):
                    target = key.removesuffix(":source") + ":translation"
                    if structures[key] != structures[target]:
                        failures.append(f"{key.removesuffix(':source')}: Go tokens, directives or output assertions changed")
        except (OSError, subprocess.CalledProcessError, ValueError) as error:
            failures.append(f"Go comment structure validation failed: {error}")
    if markdown_paths:
        try:
            result = subprocess.run(
                ["go", "run", str(ROOT / "scripts/markdown_structure.go"), *markdown_paths],
                cwd=ROOT, check=True, capture_output=True, text=True)
            structures = json.loads(result.stdout)
            for original, translated in zip(markdown_paths[::2], markdown_paths[1::2]):
                for key in ("headings", "code"):
                    if structures[original][key] != structures[translated][key]:
                        failures.append(f"{Path(translated).relative_to(CONTENT)}: Markdown {key} differ from source")
        except (OSError, subprocess.CalledProcessError, ValueError) as error:
            failures.append(f"Markdown structure validation failed: {error}")
    print(f"Checked {count} documents, including explicitly scoped partial translations.")
    for failure in failures:
        print(f"FAIL: {failure}")
    if failures:
        raise SystemExit(1)
    print("PASS: Source hashes, code, reviewed comments and prose tables, Go tokens and output assertions, reference links, unchanged unselected sections, heading IDs and Tour directives.")
    if args.base_url:
        print("PASS: Local HTTP rendering.")


if __name__ == "__main__":
    main()
