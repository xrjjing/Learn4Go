# 构建阶段
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY rabbitmq-demo ./rabbitmq-demo
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /demo ./rabbitmq-demo/cmd/demoapp

# 运行阶段
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /demo /usr/local/bin/demo
COPY rabbitmq-demo/web /app/web
COPY rabbitmq-demo/mock /app/mock
COPY rabbitmq-demo/docs /app/docs
ENV RABBITMQ_DEMO_ROOT=/app
EXPOSE 8088
ENTRYPOINT ["/usr/local/bin/demo"]
