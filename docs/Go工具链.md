# Go 工具链速览（Modules 时代）

## 环境变量
- `GOROOT`：Go 安装目录（一般无需手动修改）。  
- `GOPATH`：工作区（1.17 及以前常用，Modules 时代默认 `~/go`，仍用于缓存与 GOPATH/bin）。  
- `GOBIN`：`go install` 输出路径；未设定时为 `$GOPATH/bin`。  
- 查看：`go env`，设置：`export GOBIN=$HOME/bin`。

## 常用命令
- `go env`：查看当前环境。  
- `go mod init <module>`：创建 go.mod。  
- `go mod tidy`：解析依赖，移除未用依赖。  
- `go get <pkg>`：获取指定依赖并写入 go.mod/go.sum。  
- `go mod download`：预下载依赖（CI 可用）。  
- `go generate ./...`：按源码中的 `//go:generate` 指令生成代码/资源。

## go generate 最小示例
- 位置：`examples/basics/go_generate/`  
- 用法：`go generate ./examples/basics/go_generate/...`  
- 说明：generate 不会自动格式化或编译，需要自行 gofmt/go test。

## 编译/测试
- `go test ./...`：运行当前模块下所有包测试。  
- `go test -bench=. ./pkg/...`：运行基准测试。  
- `go vet ./...`：静态检查。  
- `go fmt ./...`：格式化。

## Modules 与 GOPATH 的关系
- Modules 时代源码不必位于 GOPATH 下；GOPATH 仍用于依赖缓存与默认 GOBIN。  
- 若需旧版 GOPATH 模式，设置 `GO111MODULE=off`（不推荐）。

## 依赖缓存
- 源码缓存：`$GOMODCACHE`（通常在 `$GOPATH/pkg/mod`）。  
- 构建缓存：`$GOCACHE`，可通过 `GOCACHE=/path` 指定（本仓库测试命令中使用）。

## 常见问题
- 依赖缺失：`go mod tidy` 或 `go get <pkg>`。  
- 私有仓库：设置 `GOPRIVATE`，或配置 git 凭据。  
- 可执行安装：`go install pkg@version`（Go1.17+ 推荐）。
