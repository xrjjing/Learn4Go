#!/bin/bash
# 全链路冒烟测试脚本：直连 TODO API +（可选）Gin/stdlib 网关

set -e

BASE_API="http://127.0.0.1:8080"
GATEWAY_API="http://127.0.0.1:8888"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}Learn4Go TODO API 全链路冒烟测试${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 工具检测
for cmd in curl jq; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo -e "${RED}❌ 缺少依赖命令: ${cmd}${NC}"
    echo "   请先安装再运行本脚本。"
    exit 1
  fi
done

echo -e "${YELLOW}[1/4] 检查 TODO API 直连健康状态...${NC}"
if ! curl -s -f "${BASE_API}/healthz" >/dev/null 2>&1; then
  echo -e "${RED}❌ 无法访问 ${BASE_API}/healthz${NC}"
  echo "   请先启动 TODO API，例如："
  echo "     go run ./cmd/todoapi"
  echo "   或使用一键脚本："
  echo "     ./start-local.sh"
  exit 1
fi
echo -e "${GREEN}✅ TODO API 健康检查通过${NC}"
echo ""

echo -e "${YELLOW}[2/4] 直连登录获取 Token（admin@example.com/admin123）...${NC}"
LOGIN_JSON=$(curl -s -X POST "${BASE_API}/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}')

TOKEN=$(echo "$LOGIN_JSON" | jq -r '.token // .access_token // .data.access_token // empty')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo -e "${RED}❌ 登录失败，未能获取 token${NC}"
  echo "   原始响应："
  echo "$LOGIN_JSON" | jq '.' || echo "$LOGIN_JSON"
  exit 1
fi

echo -e "${GREEN}✅ 登录成功，已获取访问 Token${NC}"
echo ""

echo -e "${YELLOW}[3/4] 直连访问 /v1/todos（带认证）...${NC}"
TODOS_DIRECT=$(curl -s -f "${BASE_API}/v1/todos" \
  -H "Authorization: Bearer ${TOKEN}")

if [ $? -ne 0 ]; then
  echo -e "${RED}❌ 直连 /v1/todos 访问失败${NC}"
  exit 1
fi

COUNT_DIRECT=$(echo "$TODOS_DIRECT" | jq 'length' 2>/dev/null || echo "?")
echo -e "${GREEN}✅ 直连 /v1/todos 成功，返回条目数: ${COUNT_DIRECT}${NC}"
echo ""

echo -e "${YELLOW}[4/4] 通过网关访问 /api/v1/todos（可选）...${NC}"
if ! curl -s -f "${GATEWAY_API}/health" >/dev/null 2>&1; then
  echo -e "${YELLOW}⚠️  未检测到 Gateway 运行在 ${GATEWAY_API}，跳过网关测试${NC}"
  echo "   如需验证，可在另一个终端启动："
  echo "     go run ./examples/gateway/stdlib"
  echo "   或："
  echo "     go run ./examples/gateway/gin"
  echo ""
else
  TODOS_GATEWAY=$(curl -s -f "${GATEWAY_API}/api/v1/todos" \
    -H "Authorization: Bearer ${TOKEN}")
  if [ $? -ne 0 ]; then
    echo -e "${RED}❌ 通过网关访问 /api/v1/todos 失败${NC}"
    exit 1
  fi
  COUNT_GATEWAY=$(echo "$TODOS_GATEWAY" | jq 'length' 2>/dev/null || echo "?")
  echo -e "${GREEN}✅ 通过网关 /api/v1/todos 成功，返回条目数: ${COUNT_GATEWAY}${NC}"
  echo ""
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}🎉 冒烟测试完成：直连 TODO API 已畅通${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}提示：${NC}"
echo "1. 如需验证 Docker/Nginx 链路，可在 deployments/ 目录运行 docker-compose 后手动执行："
echo "   curl http://localhost/api/v1/todos"
echo "2. 如需验证管理后台接口，请使用现有脚本："
echo "   ./test-admin-api.sh"
*** End Patch!*\
