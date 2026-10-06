#!/usr/bin/env bash
set -euo pipefail

cd "$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"

if ! command -v go >/dev/null 2>&1; then
  echo '请先安装 Go 1.26 或更新版本：https://go.dev/dl/' >&2
  exit 1
fi

# Keep the proxy setting local to this process; never change `go env -w`.
export GOPROXY="${GOPROXY:-https://goproxy.cn,https://proxy.golang.org,direct}"

export GO_WEBSITE_CONTENT_DIR="$PWD/_content"

mode=web
if [[ "${1:-}" == web || "${1:-}" == tour ]]; then
  mode="$1"
  shift
fi

if [[ "$mode" == tour ]]; then
  echo '正在启动本地 Go 语言之旅，请等待启动日志……'
  exec go run -tags=localcontent ./tour -http=127.0.0.1:3999 -openbrowser=false "$@"
fi

echo '正在准备中文文档服务，首次启动可能需要编译，请稍候……'
# PORT is reserved for upstream App Engine deployment, not local preview.
unset PORT
exec go run -tags=localcontent ./cmd/golangorg -http=127.0.0.1:6060 -content="$PWD/_content" \
  -tip=false -wiki=false -gopls=false -localdocs "$@"
