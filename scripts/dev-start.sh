#!/bin/bash
# 本地开发一键启动脚本（内存模式）
# 等价于在项目根目录执行: ./start-local.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "$PROJECT_ROOT"
./start-local.sh "$@"

