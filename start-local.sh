#!/bin/bash

# Learn4Go 本地启动脚本
# 用法: ./start-local.sh [memory|sqlite|mysql]

set -e

MODE=${1:-memory}
PROJECT_ROOT=$(cd "$(dirname "$0")" && pwd)

echo "🚀 启动 Learn4Go 本地开发环境"
echo "📂 项目目录: $PROJECT_ROOT"
echo "💾 存储模式: $MODE"
echo ""

# 检查端口占用
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        echo "⚠️  警告: 端口 $port 已被占用"
        return 1
    fi
    return 0
}

# 创建日志目录
mkdir -p logs

# 启动 TODO API
start_todoapi() {
    echo "1️⃣  启动 TODO API (端口 8080)..."

    case $MODE in
        memory)
            echo "   使用内存存储"
            nohup go run ./cmd/todoapi > logs/todoapi.log 2>&1 &
            ;;
        sqlite)
            echo "   使用 SQLite 存储"
            nohup env TODO_STORAGE=sqlite go run ./cmd/todoapi > logs/todoapi.log 2>&1 &
            ;;
        mysql)
            echo "   使用 MySQL 存储"
            echo "   确保 MySQL 已启动: docker-compose -f deployments/docker-compose.yml up -d mysql"
            nohup env TODO_STORAGE=mysql \
            TODO_DB_HOST=localhost \
            TODO_DB_PORT=3306 \
            TODO_DB_USER=root \
            TODO_DB_PASS=root \
            TODO_DB_NAME=learn4go \
            go run ./cmd/todoapi > logs/todoapi.log 2>&1 &
            ;;
        *)
            echo "❌ 未知模式: $MODE"
            echo "用法: $0 [memory|sqlite|mysql]"
            exit 1
            ;;
    esac

    TODOAPI_PID=$!
    echo "   PID: $TODOAPI_PID"
    echo "   日志: $PROJECT_ROOT/logs/todoapi.log"
    sleep 3
}

# 启动 Gateway
start_gateway() {
    echo ""
    echo "2️⃣  启动 API Gateway (端口 8888)..."

    nohup env GATEWAY_ADDR=:8888 \
    TODO_API_URL=http://localhost:8080 \
    go run ./examples/gateway/gin > logs/gateway.log 2>&1 &

    GATEWAY_PID=$!
    echo "   PID: $GATEWAY_PID"
    echo "   日志: $PROJECT_ROOT/logs/gateway.log"
    sleep 3
}

# 启动前端
start_frontend() {
    echo ""
    echo "3️⃣  启动前端服务器 (端口 8000)..."

    cd "$PROJECT_ROOT/web"
    nohup python3 -m http.server 8000 > ../logs/frontend.log 2>&1 &
    FRONTEND_PID=$!
    echo "   PID: $FRONTEND_PID"
    echo "   日志: $PROJECT_ROOT/logs/frontend.log"
    cd "$PROJECT_ROOT"
    sleep 2
}

# 健康检查
health_check() {
    echo ""
    echo "🔍 健康检查..."

    sleep 5

    # 检查 TODO API
    if curl -s http://localhost:8080/healthz > /dev/null 2>&1; then
        echo "   ✅ TODO API: http://localhost:8080"
    else
        echo "   ❌ TODO API 启动失败，查看日志: tail -f logs/todoapi.log"
    fi

    # 检查 Gateway
    if curl -s http://localhost:8888/health > /dev/null 2>&1; then
        echo "   ✅ Gateway: http://localhost:8888"
    else
        echo "   ❌ Gateway 启动失败，查看日志: tail -f logs/gateway.log"
    fi

    # 检查前端
    if curl -s http://localhost:8000 > /dev/null 2>&1; then
        echo "   ✅ Frontend: http://localhost:8000"
    else
        echo "   ❌ Frontend 启动失败，查看日志: tail -f logs/frontend.log"
    fi
}

# 显示访问信息
show_info() {
    echo ""
    echo "✨ 启动完成！"
    echo ""
    echo "📍 访问地址:"
    echo "   学习门户: http://localhost:8000/portal.html"
    echo "   项目首页: http://localhost:8000/index.html"
    echo "   项目实战: http://localhost:8000/projects.html"
    echo "   TODO API: http://localhost:8080 (根路径) 或 http://localhost:8080/todos"
    echo "   Gateway:  http://localhost:8888 (根路径) 或 http://localhost:8888/api/todos"
    echo ""
    echo "🛑 停止服务:"
    echo "   运行: ./stop-local.sh"
    echo "   或按 Ctrl+C"
    echo ""
    echo "📝 进程 ID 已保存到 logs/.pids 文件"
}

# 保存 PID
save_pids() {
    cat > "$PROJECT_ROOT/logs/.pids" <<EOF
TODOAPI_PID=$TODOAPI_PID
GATEWAY_PID=$GATEWAY_PID
FRONTEND_PID=$FRONTEND_PID
EOF
}

# 清理函数
cleanup() {
    echo ""
    echo "🛑 正在停止服务..."

    if [ -n "$TODOAPI_PID" ]; then
        kill $TODOAPI_PID 2>/dev/null || true
        echo "   已停止 TODO API"
    fi

    if [ -n "$GATEWAY_PID" ]; then
        kill $GATEWAY_PID 2>/dev/null || true
        echo "   已停止 Gateway"
    fi

    if [ -n "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
        echo "   已停止 Frontend"
    fi

    rm -f "$PROJECT_ROOT/logs/.pids"
    echo "✅ 所有服务已停止"
    exit 0
}

# 捕获退出信号
trap cleanup SIGINT SIGTERM

# 主流程
main() {
    cd "$PROJECT_ROOT"

    # 检查端口
    check_port 8080 || exit 1
    check_port 8888 || exit 1
    check_port 8000 || exit 1

    # 启动服务
    start_todoapi
    start_gateway
    start_frontend

    # 健康检查
    health_check

    # 保存 PID
    save_pids

    # 显示信息
    show_info

    # 等待
    echo "⏳ 服务运行中... (按 Ctrl+C 停止)"
    wait
}

main
