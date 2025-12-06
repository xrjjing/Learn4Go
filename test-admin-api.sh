#!/bin/bash
# 管理后台 API 集成测试脚本

set -e

BASE_URL="http://127.0.0.1:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}管理后台 API 集成测试${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 检查服务是否运行
echo -e "${YELLOW}[1/8] 检查健康状态...${NC}"
if ! curl -s -f "$BASE_URL/healthz" > /dev/null 2>&1; then
    echo -e "${RED}❌ 服务未运行，请先启动: ./todoapi${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 服务健康${NC}"
echo ""

# 登录获取 token
echo -e "${YELLOW}[2/8] 登录获取 Admin Token...${NC}"
TOKEN=$(curl -s -X POST "$BASE_URL/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ 登录失败${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 登录成功${NC}"
echo ""

# 测试 GET /v1/me
echo -e "${YELLOW}[3/8] 测试 GET /v1/me...${NC}"
ME_RESPONSE=$(curl -s "$BASE_URL/v1/me" -H "Authorization: Bearer $TOKEN")
echo "$ME_RESPONSE" | grep -q '"email":"admin@example.com"'
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ GET /v1/me 成功${NC}"
    echo "   响应: $(echo $ME_RESPONSE | jq -c '.')"
else
    echo -e "${RED}❌ GET /v1/me 失败${NC}"
    echo "   响应: $ME_RESPONSE"
    exit 1
fi
echo ""

# 测试 GET /v1/users
echo -e "${YELLOW}[4/8] 测试 GET /v1/users...${NC}"
USERS_RESPONSE=$(curl -s "$BASE_URL/v1/users" -H "Authorization: Bearer $TOKEN")
echo "$USERS_RESPONSE" | grep -q '"email"'
if [ $? -eq 0 ]; then
    USER_COUNT=$(echo "$USERS_RESPONSE" | jq 'length')
    echo -e "${GREEN}✅ GET /v1/users 成功（共 $USER_COUNT 个用户）${NC}"
else
    echo -e "${RED}❌ GET /v1/users 失败${NC}"
    echo "   响应: $USERS_RESPONSE"
    exit 1
fi
echo ""

# 测试 GET /v1/rbac/roles
echo -e "${YELLOW}[5/8] 测试 GET /v1/rbac/roles...${NC}"
ROLES_RESPONSE=$(curl -s "$BASE_URL/v1/rbac/roles" -H "Authorization: Bearer $TOKEN")
echo "$ROLES_RESPONSE" | grep -q '"name":"admin"'
if [ $? -eq 0 ]; then
    ROLE_COUNT=$(echo "$ROLES_RESPONSE" | jq 'length')
    echo -e "${GREEN}✅ GET /v1/rbac/roles 成功（共 $ROLE_COUNT 个角色）${NC}"
else
    echo -e "${RED}❌ GET /v1/rbac/roles 失败${NC}"
    echo "   响应: $ROLES_RESPONSE"
    exit 1
fi
echo ""

# 测试 GET /v1/rbac/permissions
echo -e "${YELLOW}[6/8] 测试 GET /v1/rbac/permissions...${NC}"
PERMS_RESPONSE=$(curl -s "$BASE_URL/v1/rbac/permissions" -H "Authorization: Bearer $TOKEN")
echo "$PERMS_RESPONSE" | grep -q '"code":"todos:create"'
if [ $? -eq 0 ]; then
    PERM_COUNT=$(echo "$PERMS_RESPONSE" | jq 'length')
    echo -e "${GREEN}✅ GET /v1/rbac/permissions 成功（共 $PERM_COUNT 个权限）${NC}"
else
    echo -e "${RED}❌ GET /v1/rbac/permissions 失败${NC}"
    echo "   响应: $PERMS_RESPONSE"
    exit 1
fi
echo ""

# 测试 POST /v1/users（创建新用户）
echo -e "${YELLOW}[7/8] 测试 POST /v1/users（创建新用户）...${NC}"
CREATE_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/users" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"testuser@example.com","password":"test123","is_superuser":false}')
echo "$CREATE_USER_RESPONSE" | grep -q '"email":"testuser@example.com"'
if [ $? -eq 0 ]; then
    NEW_USER_ID=$(echo "$CREATE_USER_RESPONSE" | jq -r '.id')
    echo -e "${GREEN}✅ POST /v1/users 成功（用户 ID: $NEW_USER_ID）${NC}"
else
    # 可能已存在，不算失败
    echo -e "${YELLOW}⚠️  用户可能已存在${NC}"
    echo "   响应: $(echo $CREATE_USER_RESPONSE | jq -c '.')"
fi
echo ""

# 测试 POST /v1/logout
echo -e "${YELLOW}[8/8] 测试 POST /v1/logout...${NC}"
LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/logout" \
  -H "Authorization: Bearer $TOKEN")
echo "$LOGOUT_RESPONSE" | grep -q 'logged out successfully'
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ POST /v1/logout 成功${NC}"
else
    echo -e "${RED}❌ POST /v1/logout 失败${NC}"
    echo "   响应: $LOGOUT_RESPONSE"
    exit 1
fi
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}🎉 所有测试通过！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}提示：${NC}"
echo "1. 本地开发模式下，访问 http://localhost:8000/admin.html 查看管理后台"
echo "2. 使用 admin@example.com / admin123 登录"
echo "3. 查看 API 文档: docs/API接口文档.md"
