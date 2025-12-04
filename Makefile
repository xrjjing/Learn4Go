GOFILES := $(shell find . -name '*.go' -not -path './.git/*')

.PHONY: fmt vet test lint tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || echo "golangci-lint 未安装，跳过（可自行安装后运行 make lint）"

tidy:
	go mod tidy
