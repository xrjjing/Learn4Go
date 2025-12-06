#!/bin/bash
# Docker 部署：停止全部服务（保留数据卷）
# 等价于在 deployments 目录执行: docker-compose down

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "$PROJECT_ROOT/deployments"
docker-compose down "$@"

