#!/bin/bash

# Learn4Go 临时文件清理脚本

PROJECT_ROOT=$(cd "$(dirname "$0")" && pwd)

echo "🧹 清理 Learn4Go 临时文件..."

# 停止所有服务
if [ -f "$PROJECT_ROOT/stop-local.sh" ]; then
    echo "   停止运行中的服务..."
    ./stop-local.sh 2>/dev/null || true
fi

# 清理日志目录
if [ -d "$PROJECT_ROOT/logs" ]; then
    echo "   清理日志目录..."
    rm -rf "$PROJECT_ROOT/logs"
fi

# 清理编译的二进制文件
echo "   清理编译产物..."
rm -f "$PROJECT_ROOT/client"
rm -f "$PROJECT_ROOT/server"
rm -f "$PROJECT_ROOT/gin"
rm -f "$PROJECT_ROOT/gateway"
rm -f "$PROJECT_ROOT/todoapi"

# 清理数据库文件
echo "   清理数据库文件..."
rm -f "$PROJECT_ROOT"/*.db
rm -f "$PROJECT_ROOT"/*.sqlite
rm -f "$PROJECT_ROOT"/*.sqlite3

# 清理其他临时文件
echo "   清理其他临时文件..."
rm -f "$PROJECT_ROOT/nohup.out"
rm -f "$PROJECT_ROOT"/*.log

echo ""
echo "✅ 清理完成！"
echo ""
echo "📊 当前状态:"
echo "   日志目录: $([ -d logs ] && echo '存在' || echo '已清理')"
echo "   二进制文件: $(ls -1 client server gin gateway todoapi 2>/dev/null | wc -l | tr -d ' ') 个"
echo "   数据库文件: $(ls -1 *.db *.sqlite *.sqlite3 2>/dev/null | wc -l | tr -d ' ') 个"
echo ""
echo "💡 提示: 运行 ./start-local.sh 重新启动服务"
