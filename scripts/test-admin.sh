#!/bin/bash
# 管理后台 API 集成测试包装脚本
# 等价于在项目根目录执行: ./test-admin-api.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "$PROJECT_ROOT"
./test-admin-api.sh "$@"

