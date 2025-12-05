#!/bin/bash

# Learn4Go 本地服务停止脚本

PROJECT_ROOT=$(cd "$(dirname "$0")" && pwd)
PID_FILE="$PROJECT_ROOT/logs/.pids"

echo "🛑 停止 Learn4Go 本地服务..."

if [ ! -f "$PID_FILE" ]; then
    echo "⚠️  未找到 PID 文件，尝试通过端口查找进程..."

    # 通过端口查找并停止
    for port in 8080 8888 8000; do
        pid=$(lsof -ti:$port 2>/dev/null)
        if [ -n "$pid" ]; then
            echo "   停止端口 $port 上的进程 (PID: $pid)"
            kill $pid 2>/dev/null || true
        fi
    done
else
    # 从文件读取 PID
    source "$PID_FILE"

    if [ -n "$TODOAPI_PID" ]; then
        kill $TODOAPI_PID 2>/dev/null && echo "   ✅ 已停止 TODO API (PID: $TODOAPI_PID)" || echo "   ⚠️  TODO API 进程不存在"
    fi

    if [ -n "$GATEWAY_PID" ]; then
        kill $GATEWAY_PID 2>/dev/null && echo "   ✅ 已停止 Gateway (PID: $GATEWAY_PID)" || echo "   ⚠️  Gateway 进程不存在"
    fi

    if [ -n "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null && echo "   ✅ 已停止 Frontend (PID: $FRONTEND_PID)" || echo "   ⚠️  Frontend 进程不存在"
    fi

    rm -f "$PID_FILE"
fi

echo ""
echo "✅ 所有服务已停止"
