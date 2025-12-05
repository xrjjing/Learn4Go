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

.PHONY: tinygee-test tinygee-bench
tinygee-test:
	GOCACHE=$(PWD)/.gocache go test ./tinygee/...

tinygee-bench:
	GOCACHE=$(PWD)/.gocache go test -bench=. ./tinygee/...

.PHONY: tinygee-examples
tinygee-examples:
	@echo "示例入口:"
	@echo "  go run ./examples/tinygee/day1"
	@echo "  go run ./examples/tinygee/day3"
	@echo "  go run ./examples/tinygee/day5"
	@echo "  go run ./examples/tinygee/day6"
	@echo "  go run ./examples/tinygee/day7"
	@echo "  go run ./examples/tinygee/day9auth"
	@echo "  go run ./examples/tinygee/day9rl"
.PHONY: demo-java
demo-java:
	@echo "按需运行（部分依赖 httpbin 需联网）："
	@echo "go run ./examples/java_compare/urlencode"
	@echo "go run ./examples/java_compare/bufio_wordcount"
	@echo "go run ./examples/java_compare/concurrency"
	@echo "go run ./examples/java_compare/context_timeout"
	@echo "go run ./examples/java_compare/ticker_rate_limit"
	@echo "go run ./examples/java_compare/httptrace"
	@echo "go run ./examples/java_compare/http_middleware"
	@echo "go run ./examples/java_compare/syncmap"
	@echo "go run ./examples/java_compare/pprof_server"
