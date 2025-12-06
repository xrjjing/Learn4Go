#!/bin/bash
# 本地开发一键停止脚本
# 等价于在项目根目录执行: ./stop-local.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "$PROJECT_ROOT"
./stop-local.sh "$@"

