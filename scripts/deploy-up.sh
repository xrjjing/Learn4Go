#!/bin/bash
# Docker 部署：启动全部服务（前端 + Gateway + TODO API 等）
# 等价于在 deployments 目录执行: docker-compose up -d

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "$PROJECT_ROOT/deployments"
docker-compose up -d "$@"

