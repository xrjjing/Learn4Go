# 贡献指南

感谢您有兴趣为 Learn4Go 项目做出贡献！

## 开发环境

### 前置要求

- Go 1.21+
- Node.js 18+ (前端开发)
- Docker (可选，用于数据库)

### 快速开始

```bash
# 克隆项目
git clone https://github.com/xrjjing/Learn4Go.git
cd Learn4Go

# 安装依赖
go mod download

# 运行测试
go test ./...

# 启动服务
go run ./cmd/todoapi
```

## 代码规范

### Go 代码

- 遵循 [Effective Go](https://go.dev/doc/effective_go)
- 使用 `gofmt` 格式化代码
- 导出函数必须有注释
- 错误处理：不忽略错误，使用 `%w` 包装

```bash
# 格式化
gofmt -w .

# 静态检查
go vet ./...
```

### 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

类型：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具

示例：
```
feat(api): add pagination to /v1/todos endpoint

- Add page and page_size query parameters
- Return total count in response header
```

## 提交 PR

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feat/your-feature`)
3. 提交更改 (`git commit -m 'feat: add new feature'`)
4. 推送分支 (`git push origin feat/your-feature`)
5. 创建 Pull Request

### PR 检查清单

- [ ] 代码通过 `go build` 和 `go test`
- [ ] 新功能有对应测试
- [ ] 更新相关文档
- [ ] 提交信息符合规范

## 项目结构

```
Learn4Go/
├── cmd/              # 可执行入口
│   └── todoapi/      # TODO API 服务
├── internal/         # 内部包
│   └── todo/         # TODO 业务逻辑
├── web/              # 前端资源
├── docs/             # 文档
└── examples/         # 示例代码
```

## 问题反馈

- 使用 GitHub Issues 报告 Bug
- 描述清楚复现步骤
- 提供错误日志和环境信息

## 许可证

MIT License
