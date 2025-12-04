# gRPC 示例的 Makefile
#
# 使用方法：
#   make proto  - 从 .proto 文件生成 Go 代码
#   make clean  - 清理生成的代码
#   make help   - 显示帮助信息

PROTOC ?= protoc
PROTO_DIR := proto
OUT_DIR := .

# 生成 Go 代码
# --go_out: 生成消息结构体
# --go-grpc_out: 生成 gRPC 服务接口
proto:
	@echo "生成 protobuf 代码..."
	$(PROTOC) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_DIR)/hello.proto
	@echo "生成完成！文件位于: $(OUT_DIR)/$(PROTO_DIR)/hellopb/"

# 清理生成的文件
clean:
	rm -rf $(PROTO_DIR)/hellopb

# 显示帮助
help:
	@echo "gRPC 示例 Makefile"
	@echo ""
	@echo "目标:"
	@echo "  proto  - 生成 protobuf Go 代码"
	@echo "  clean  - 清理生成的代码"
	@echo "  help   - 显示此帮助信息"
	@echo ""
	@echo "前置要求:"
	@echo "  1. 安装 protoc 编译器"
	@echo "  2. 安装 Go 插件:"
	@echo "     go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
	@echo "     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"

.PHONY: proto clean help
