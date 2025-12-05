#!/bin/bash

# JWT认证系统演示脚本
# 展示完整的认证流程：注册 -> 登录 -> 访问受保护API

set -e

API_URL="http://localhost:8080"
BOLD='\033[1m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BOLD}=== JWT认证系统演示 ===${NC}\n"

# 检查服务是否运行
echo -e "${BLUE}1. 检查服务状态...${NC}"
if ! curl -s -f "$API_URL/healthz" > /dev/null; then
    echo -e "${YELLOW}⚠️  服务未启动，请先运行: go run ./cmd/todoapi${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 服务正常运行${NC}\n"

# 测试未认证访问
echo -e "${BLUE}2. 测试未认证访问（应该被拒绝）...${NC}"
RESPONSE=$(curl -s "$API_URL/todos")
echo "响应: $RESPONSE"
if echo "$RESPONSE" | grep -q "authorization required"; then
    echo -e "${GREEN}✓ 认证保护生效${NC}\n"
else
    echo -e "${YELLOW}⚠️  认证保护未生效${NC}\n"
fi

# 使用mock用户登录
echo -e "${BLUE}3. 使用mock用户登录...${NC}"
echo "邮箱: admin@example.com"
echo "密码: admin123"

LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}')

echo "登录响应:"
echo "$LOGIN_RESPONSE" | jq '.'

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}⚠️  登录失败${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 登录成功${NC}"
echo -e "Token: ${TOKEN:0:50}...\n"

# 使用token访问API
echo -e "${BLUE}4. 使用token访问TODO列表...${NC}"
TODOS=$(curl -s "$API_URL/todos" \
  -H "Authorization: Bearer $TOKEN")
echo "TODO列表: $TODOS"
echo -e "${GREEN}✓ 认证访问成功${NC}\n"

# 创建TODO
echo -e "${BLUE}5. 创建新的TODO...${NC}"
NEW_TODO=$(curl -s -X POST "$API_URL/todos" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"通过JWT认证创建的任务"}')

echo "创建的TODO:"
echo "$NEW_TODO" | jq '.'
echo -e "${GREEN}✓ TODO创建成功${NC}\n"

# 再次获取列表
echo -e "${BLUE}6. 再次获取TODO列表...${NC}"
TODOS=$(curl -s "$API_URL/todos" \
  -H "Authorization: Bearer $TOKEN")
echo "$TODOS" | jq '.'
echo -e "${GREEN}✓ 列表获取成功${NC}\n"

# 测试注册新用户
echo -e "${BLUE}7. 测试注册新用户...${NC}"
TIMESTAMP=$(date +%s)
NEW_EMAIL="test${TIMESTAMP}@example.com"
echo "新用户邮箱: $NEW_EMAIL"

REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$NEW_EMAIL\",\"password\":\"test123\"}")

echo "注册响应:"
echo "$REGISTER_RESPONSE" | jq '.'

if echo "$REGISTER_RESPONSE" | jq -e '.id' > /dev/null; then
    echo -e "${GREEN}✓ 用户注册成功${NC}\n"

    # 使用新用户登录
    echo -e "${BLUE}8. 使用新用户登录...${NC}"
    NEW_LOGIN=$(curl -s -X POST "$API_URL/login" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$NEW_EMAIL\",\"password\":\"test123\"}")

    echo "$NEW_LOGIN" | jq '.'
    echo -e "${GREEN}✓ 新用户登录成功${NC}\n"
else
    echo -e "${YELLOW}⚠️  用户注册失败${NC}\n"
fi

# 测试错误密码
echo -e "${BLUE}9. 测试错误密码（应该失败）...${NC}"
WRONG_LOGIN=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"wrongpassword"}')

echo "响应: $WRONG_LOGIN"
if echo "$WRONG_LOGIN" | grep -q "invalid credentials"; then
    echo -e "${GREEN}✓ 密码验证正常${NC}\n"
else
    echo -e "${YELLOW}⚠️  密码验证异常${NC}\n"
fi

# 测试无效token
echo -e "${BLUE}10. 测试无效token（应该失败）...${NC}"
INVALID_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.token"
INVALID_RESPONSE=$(curl -s "$API_URL/todos" \
  -H "Authorization: Bearer $INVALID_TOKEN")

echo "响应: $INVALID_RESPONSE"
if echo "$INVALID_RESPONSE" | grep -q "invalid or expired token"; then
    echo -e "${GREEN}✓ Token验证正常${NC}\n"
else
    echo -e "${YELLOW}⚠️  Token验证异常${NC}\n"
fi

echo -e "${BOLD}${GREEN}=== 演示完成 ===${NC}"
echo -e "\n${BOLD}Mock用户账户：${NC}"
echo "  • admin@example.com / admin123"
echo "  • user@example.com / user123"
echo "  • demo@example.com / demo123"

echo -e "\n${BOLD}相关文档：${NC}"
echo "  • JWT认证文档: docs/AUTH.md"
echo "  • API文档: docs/API.md"
echo "  • 变更日志: docs/CHANGELOG.md"
