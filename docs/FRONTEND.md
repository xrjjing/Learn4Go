# 前端使用指南

本文档详细介绍 Learn4Go 项目前端页面的使用方法。

## 📂 前端文件结构

```
web/
├── portal.html          # 学习门户主页（推荐入口）
├── config.js            # 前端配置文件
├── index.html           # 项目首页
├── projects.html        # 项目列表页
├── admin.html           # 管理页面
├── login.html           # 旧版登录页面（RBAC 示例）
├── todo-login.html      # TODO API 登录 + Refresh 示例
├── log-detective.html   # 日志分析工具
├── common.js            # 公共 JavaScript
├── auth-helper.js       # Token 管理与自动刷新
├── mock-data.js         # Mock 数据
├── mock-api.js          # Mock API
└── README.md            # 前端说明
```

## 🚀 快速开始

### 方式一：Docker 部署（推荐）

```bash
# 1. 启动所有服务
cd deployments
docker-compose up -d

# 2. 访问前端
open http://localhost
```

所有服务会自动启动，前端通过 Nginx 提供服务。

### 方式二：本地开发

```bash
# 1. 启动后端服务（参考 README.md）
go run ./cmd/todoapi
go run ./examples/gateway/gin

# 2. 启动前端服务器
cd web
python3 -m http.server 8000

# 3. 访问前端
open http://localhost:8000/portal.html
```

## 🎨 页面功能详解

### 1. 学习门户 (portal.html)

**访问地址**: http://localhost/portal.html

这是项目的主要入口页面，提供完整的学习体验。

#### 功能特性

##### 1.1 服务状态监控

页面右上角实时显示后端服务状态：

- **TodoAPI**: TODO REST API 服务状态
- **Gateway**: API 网关服务状态

状态指示：
- 🟢 绿色：服务在线
- 🔴 红色：服务离线
- ⚪ 灰色：检测中

##### 1.2 学习进度追踪

顶部进度条显示整体学习进度：

- 自动保存到浏览器 LocalStorage
- 显示已完成章节数 / 总章节数
- 支持重置进度（点击"重置进度"按钮）

##### 1.3 章节导航

左侧边栏分为三大模块：

**基础章节**（10 个）
1. 快速开始与基本语法
2. 变量、常量、类型
3. 流程控制
4. 函数与错误处理
5. 数组、切片、Map
6. 结构体、方法、组合
7. 接口与多态
8. 并发：goroutine 与 channel
9. 包与模块管理
10. 测试与基准

**开发环境与框架**（4 个）
1. 工具链与依赖管理
2. Web 框架对比
3. gRPC 与微服务入门
4. 日志、配置与中间件

**项目实战**（6 个）
1. 批量重命名 CLI
2. TODO REST API
3. 日志 ZIP 归档
4. Worker Pool 并发
5. API 网关（标准库）
6. API 网关（Gin）

##### 1.4 可运行示例

页面底部列出所有可运行的示例代码：

- **示例名称**: 如 "Hello World"、"FizzBuzz"
- **文件路径**: 源代码位置
- **运行命令**: 可一键复制的命令
- **操作按钮**:
  - 📋 复制：复制运行命令到剪贴板
  - ▶ 运行：显示运行提示

##### 1.5 章节学习流程

1. 点击左侧章节名称
2. 查看章节内容和示例代码
3. 复制运行命令到终端执行
4. 完成后点击"标记为已完成"
5. 进度自动保存

### 2. 项目首页 (index.html)

**访问地址**: http://localhost/index.html

项目的欢迎页面，展示项目概览和快速链接。

### 3. 项目列表 (projects.html)

**访问地址**: http://localhost/projects.html

展示所有实战项目的详细信息和运行方法。

### 4. 日志分析工具 (log-detective.html)

**访问地址**: http://localhost/log-detective.html

一个简单的日志分析工具，用于演示前后端交互。

## ⚙️ 配置说明

### config.js 配置文件

前端配置文件会自动检测运行环境：

```javascript
// Docker 部署模式
{
    todoApiBaseUrl: '/api/todos',
    gatewayUrl: '/api',
    enableMock: false
}

// 本地开发模式
{
    todoApiBaseUrl: 'http://127.0.0.1:8080',
    gatewayUrl: 'http://127.0.0.1:8888',
    enableMock: true
}
```

#### 环境检测逻辑

```javascript
const isDocker = window.location.hostname !== 'localhost'
                 && window.location.hostname !== '127.0.0.1';
```

- **Docker 模式**: 使用相对路径，Nginx 代理到后端
- **本地模式**: 使用绝对地址，直连后端服务

### 修改配置

如需修改配置，编辑 `web/config.js`：

```javascript
window.AppConfig = {
    todoApiBaseUrl: 'http://your-api-url',
    gatewayUrl: 'http://your-gateway-url',
    enableMock: false
};
```

## 🔧 开发调试

### 查看浏览器控制台

按 F12 打开开发者工具，查看：

- **Console**: JavaScript 日志和错误
- **Network**: API 请求和响应
- **Application > Local Storage**: 学习进度数据

### 清除学习进度

在浏览器控制台执行：

```javascript
localStorage.removeItem('learn4goProgress');
location.reload();
```

或点击页面上的"重置进度"按钮。

### Mock 模式

本地开发时，如果后端服务未启动，前端会自动使用 Mock 数据：

- `mock-data.js`: Mock 数据定义
- `mock-api.js`: Mock API 实现

## 🎯 使用场景

### 场景一：学习 Go 语言

1. 访问 http://localhost/portal.html
2. 从"基础章节"开始学习
3. 按顺序完成每个章节
4. 运行示例代码验证理解
5. 标记完成并查看进度

### 场景二：测试 TODO API

1. 确保 TODO API 服务运行
2. 查看服务状态指示器（应为绿色）
3. 使用浏览器或 curl 测试 API：

```bash
# 获取所有 TODO
curl http://localhost/api/todos

# 创建 TODO
curl -X POST http://localhost/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"学习 Go 语言"}'

# 更新 TODO
curl -X PUT http://localhost/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"done":true}'

# 删除 TODO
curl -X DELETE http://localhost/api/todos/1
```

### 场景三：演示微服务架构

1. 启动完整的 Docker 环境
2. 访问前端门户
3. 观察服务状态监控
4. 通过 API 网关访问后端服务
5. 查看 MinIO 控制台（http://localhost:9001）

## 📱 响应式设计

前端页面支持多种设备：

- **桌面**: 完整的侧边栏和内容区域
- **平板**: 自适应布局
- **手机**: 折叠侧边栏，垂直布局

## 🎨 UI 组件

### Bootstrap 5 组件

- **导航栏**: navbar
- **侧边栏**: sidebar
- **卡片**: card
- **按钮**: btn
- **徽章**: badge
- **进度条**: progress

### 自定义样式

- **字体**: Inter（正文）+ Fira Code（代码）
- **配色**:
  - 主色：#0d6efd（蓝色）
  - 成功：#198754（绿色）
  - 危险：#dc3545（红色）
  - 背景：#f8f9fa（浅灰）

## 🔍 故障排查

### 问题：服务状态显示离线

**解决方案**：

1. 检查后端服务是否启动：
   ```bash
   docker-compose ps
   # 或
   curl http://localhost:8080/healthz
   ```

2. 检查浏览器控制台是否有 CORS 错误

3. 确认 Nginx 配置正确（Docker 模式）

### 问题：学习进度丢失

**解决方案**：

1. 检查浏览器是否禁用了 LocalStorage
2. 不要使用隐私/无痕模式
3. 检查浏览器控制台错误

### 问题：页面样式错乱

**解决方案**：

1. 清除浏览器缓存（Ctrl+Shift+R）
2. 检查 CDN 资源是否加载成功
3. 检查网络连接

## 📚 相关文档

- [README.md](../README.md) - 项目总览
- [部署指南](../deployments/README.md) - Docker 部署
- [API 文档](API.md) - TODO API 接口
- [项目计划](../plan.md) - 后续开发计划

## 💡 最佳实践

1. **使用 Docker 部署**: 最简单的方式，一键启动所有服务
2. **按顺序学习**: 从基础章节开始，循序渐进
3. **动手实践**: 运行每个示例代码，加深理解
4. **记录进度**: 及时标记完成的章节
5. **查看源码**: 学习前端代码实现，理解前后端交互

## 🔐 登录与鉴权演示
- `login.html`：面向 `rbac_auth_service`（Python 版） 的旧示例。  
- `todo-login.html`：Go TODO API 的新示例，使用 `auth-helper.js` 自动刷新 token。收到 401 时自动调用 `POST /refresh` 并重试请求。
- Mock 说明：`login.html` 支持 `?mock=true` 或 `localStorage.setItem('mockApi','true')` 启用 mock-api；`todo-login.html` 需真实 TODO 后端（可用 `./start-local.sh` 启动）。
- 门户入口：`portal.html` 顶部导航和首页卡片已区分 “真实后端” 与 “Mock 演示” 按钮。

目录结构小结：
- `web/todo-login.html` + `web/auth-helper.js` → Go TODO API 真实链路演示  
- `web/login.html` + `web/mock-api.js` + `web/mock-data.js` → Mock / Python 版 RBAC 演示  
- `web/config.js` → 环境基址与 mock 开关
- 门户“可运行示例”列表会标注：`🌐` 需外网/httpbin，`本地` 本地可跑，`Mock` 可直接在前端演示。

快速体验：
```bash
go run ./cmd/todoapi
cd web && python3 -m http.server 8000
# 浏览器访问 http://localhost:8000/todo-login.html
```

前端入口：门户导航新增“登录演示”按钮（`portal.html` 顶部导航 & 卡片）。  
提示：默认账户 `admin@example.com/admin123`；若后端未启动请执行 `./start-local.sh`。

## 🎉 开始学习

现在就访问 http://localhost/portal.html 开始你的 Go 语言学习之旅吧！
