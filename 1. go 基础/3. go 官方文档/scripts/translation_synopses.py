#!/usr/bin/env python3
"""Refresh Chinese package synopses and snapshot subdirectory tables."""
from pathlib import Path
import re, json, html, os

ROOT = Path(__file__).resolve().parents[1]
os.chdir(ROOT)

def units(text):
    # Go's generated documentation omits optional paragraph closing tags.
    return re.finditer(r'<(?P<tag>p)\b[^>]*>(?P<body>.*?)(?=</?(?:p|h[1-6]|li|dt|dd|pre|ul|ol|dl|div|article)\b|\Z)', text, re.S)

synopses={}
for path in Path('_content/pkg').rglob('index.html'):
 s=path.read_text();a=s.find('id="pkg-overview"');b=s.find('id="pkg-index"',a)
 if a<0:continue
 paras=[u['body'] for u in units(s[a:b]) if u['tag']=='p']
 if not paras:continue
 text=html.unescape(re.sub('<[^>]+>','',paras[0]));text=' '.join(text.split());text=text.split('。')[0].rstrip('：:；;')+'。'
 if re.search('[\u4e00-\u9fff]',text):synopses[str(path.parent.relative_to('_content/pkg'))]=text
for path in Path('_content/cmd').glob('*/index.html'):
 paras=[u['body'] for u in units(path.read_text()) if u['tag']=='p']
 for para in paras:
  text=html.unescape(re.sub('<[^>]+>','',para));text=' '.join(text.split())
  if re.search('[\u4e00-\u9fff]',text) and not text.startswith(('中文译文','译注','本页')):
   synopses['cmd/'+path.parent.name]=text.split('。')[0].rstrip('：:；;')+'。';break
synopses['cmd/go']='go 是用于管理 Go 源码的工具。'
synopses['arena']='arena 实验包支持批量分配内存，并手动一次性释放；默认文档视图仅提供概述。'
synopses.update({
 'crypto/boring':'boring 包提供仅在使用 Go+BoringCrypto 构建时可用的函数。（特殊构建文档保留英文。）',
 'crypto/tls/fipsonly':'fipsonly 包将所有 TLS 配置限制为经 FIPS 批准的设置。（特殊构建文档保留英文。）',
 'runtime/cgo':'cgo 包为 cgo 工具生成的代码提供运行时支持。（启用 cgo 的构建视图保留英文。）',
 'runtime/secret':'secret 包提供清零辅助函数，用于清除用户程序通常不可见的内存，以实现前向保密。（实验构建文档保留英文。）',
 'simd':'simd 包实现可移植且不依赖向量大小的 SIMD 类型，以及操作这些类型的函数和方法。（实验构建文档保留英文。）',
 'simd/archsimd':'archsimd 包提供特定架构的 SIMD 操作。（实验构建文档保留英文。）',
 'syscall/js':'js 包提供在 js/wasm 架构下访问 WebAssembly 宿主环境的能力。（其他平台文档保留英文。）',
})
Path('_content/pkg-synopses.json').write_text(json.dumps(synopses,ensure_ascii=False,indent=2,sort_keys=True)+'\n')
# Snapshot subdirectory rows use the synopsis for the referenced package.
changes=0
for name in Path('translations/completed.txt').read_text().splitlines():
 if not name.startswith('pkg/'):continue
 path=Path('_content',name);s=path.read_text();prefix=str(path.parent.relative_to('_content/pkg'))
 def row(m):
  global changes
  a=re.search(r'href="([^"]+)"',m[0]);desc=re.search(r'(<td class="pkg-synopsis">)(.*?)(</td>)',m[0],re.S)
  if not a or not desc:return m[0]
  url=a[1].split('?')[0].rstrip('/');key=url.removeprefix('/pkg/') if url.startswith('/pkg/') else prefix+'/'+url
  if key not in synopses:return m[0]
  result=m[0][:desc.start(2)]+html.escape(synopses[key])+m[0][desc.end(2):]
  changes+=result!=m[0];return result
 s=re.sub(r'<tr>.*?</tr>',row,s,flags=re.S);path.write_text(s)
print('Chinese synopses:',len(synopses),'updated directory rows:',changes)
