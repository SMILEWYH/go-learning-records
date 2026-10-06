#!/usr/bin/env python3
"""Build a reviewable translation inventory; completed.txt is explicitly maintained."""

import hashlib
import json
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / "translations/en"


def title(text, fallback):
    if text.lstrip().startswith("<codewalk") or "<codewalk title=" in text:
        return re.search(r'<codewalk title="([^"]+)"', text)[1]
    match = re.search(r'"[Tt]itle"\s*:\s*"([^"]+)"|^(?:title|linkTitle):\s*(.+)$', text, re.M)
    return (match[1] or match[2]).strip('"\' ') if match else fallback


completed = set((ROOT / "translations/completed.txt").read_text().splitlines())
partial = json.loads((ROOT / "translations/partial.json").read_text())
references = json.loads((ROOT / "translations/reference-sources.json").read_text())
entries = []
for original in sorted(SOURCE.rglob("*")):
    if not original.is_file():
        continue
    path = original.relative_to(SOURCE).as_posix()
    if path.startswith("tour/") and not (original.parent == SOURCE / "tour" and original.suffix == ".article"):
        continue
    source = original.read_text(encoding="utf-8")
    current = (ROOT / "_content" / path).read_text(encoding="utf-8")
    status = "translated" if path in completed else "pending"
    if path in partial:
        status = "partial"
    if (re.fullmatch(r"doc/go1(?:\.\d+)?\.md", path)
            or path.startswith("doc/devel/") or path == "doc/next.md"):
        status = "excluded-history"
    if re.search(r'"Redirect"\s*:|^redirect:', source, re.M):
        status = "redirect"
    if path in {"doc/articles/wiki/edit.html", "doc/articles/wiki/view.html"}:
        status = "code-example"
    if path.endswith(".article"):
        label = current.splitlines()[0] if status == "translated" else source.splitlines()[0]
        url = "/tour/" + original.stem + "/1"
    else:
        label = title(current if status in {"translated", "partial"} else source, path)
        url = "/" + str(Path(path).with_suffix(""))
        if original.stem == "index":
            url = "/" + str(Path(path).parent) + "/"
    entries.append({"path": path, "title": label, "url": url, "status": status,
                    "source_sha256": hashlib.sha256(source.encode()).hexdigest()})
    if path in partial:
        entries[-1]["coverage"] = partial[path]

manifest = {"language": "zh-CN", "scope": "current technical documentation, linux/amd64 public standard library, command manuals and website pages; blogs, historical talks and release notes excluded",
            "reference_snapshot": references, "documents": entries}
(ROOT / "translations/manifest.json").write_text(json.dumps(manifest, ensure_ascii=False, indent=2) + "\n")
counts = {status: sum(e["status"] == status for e in entries)
          for status in ("translated", "partial", "pending", "excluded-history", "redirect", "code-example")}
labels = {"translated": "已翻译", "partial": "部分章节已翻译", "pending": "待翻译", "excluded-history": "版本说明（保留英文）",
          "redirect": "重定向入口", "code-example": "原始代码示例"}
lines = ["# 中文翻译清单", "", "；".join(f"{labels[key]}：{value}" for key, value in counts.items()), "",
         "由 `python3 scripts/translation_status.py` 根据明确登记的完成清单生成。",
         "`completed.txt` 登记已完成的文档与页面条目；外链案例仅包含本站标题和简介，不代表外部文章全文已译。历史版本说明、博客和历史演讲保留英文。",
         "语言规范与 go 命令手册已全文翻译；`partial.json` 仅用于登记未来尚未译完的章节范围。",
         "标准库、命令及参考文档的原文版本见 `reference-sources.json`；程序逻辑和输出保留原文，已译标准库的自然语言注释采用中文。",
         "", "| 状态 | 文档 | 源文件 |", "| --- | --- | --- |"]
for entry in entries:
    lines.append(f'| {labels[entry["status"]]} | {entry["title"]} | {entry["path"]} |')
(ROOT / "translations/STATUS.md").write_text("\n".join(lines) + "\n")
coverage = [f'已完成 **{counts["translated"]} 项文档、包参考、课程与网站页面条目**的翻译（程序逻辑与输出保留原文，已译标准库的字段和示例注释采用中文）。']
if counts["pending"]:
    coverage.append(f'还有 {counts["pending"]} 项正在处理，以下清单明确标出状态。')
else:
    coverage.append("已完成清单中登记的翻译范围；Go 语言之旅按 7 组课程计数；网站案例卡片也按独立源文件计数。公开标准库采用 linux/amd64 快照，内部实现包、其他平台视图和外部网站不在此范围内。")
if counts["partial"]:
    coverage.append(f'另有 **{counts["partial"]} 份参考文档的重点章节**已翻译；其余章节保留英文，不计入全文完成数量。参阅[中文参考文档阅读入口](/doc/reference-zh)。')
coverage += [f'版本说明和早期更新记录共 {counts["excluded-history"]} 份，按本次约定保留英文。',
             "英文原文保存在项目的 `translations/en/`，逐篇状态与原文哈希见 `translations/manifest.json`。标准库与命令参考快照为 Go 1.27.1、linux/amd64。"]
table = ["| 文档或课程 | 状态 |", "| --- | --- |"]
for entry in entries:
    if entry["status"] in {"translated", "partial", "pending"}:
        table.append(f'| [{entry["title"]}]({entry["url"]}) | {labels[entry["status"]]} |')
page = ROOT / "_content/doc/translation.md"
page.write_text(re.sub(r"(?s)(<!-- translation-coverage:start -->).*?(<!-- translation-coverage:end -->)",
                       lambda match: match[1] + "\n\n" + "\n\n".join(coverage) + "\n\n" + "\n".join(table) + "\n\n" + match[2],
                       page.read_text()))
print(counts)
