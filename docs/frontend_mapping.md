# 前端与 Go 后端映射说明

## 静态资源
- 目录：`web/`
- mock 开关：URL `?mock=true` 或 `localStorage.setItem('mockApi','true')`

## 已对接的 Go 服务
- TODO API：`cmd/todoapi`，端口 `8080`
  - `GET /todos` 列表
  - `POST /todos` 创建
  - `PUT /todos/{id}` 更新完成状态
  - `DELETE /todos/{id}` 删除
- 配置：`web/config.js` 中 `todoApiBaseUrl`

## 预留（待实现）
- 网关/认证/日志服务：保持与 Python 仓库路径一致
  - `apiBaseUrl`（未来 gateway）
  - `logApiBaseUrl`（审计/日志服务）
- 前端页面：`web/login.html`、`web/admin.html`、`web/log-detective.html`，将逐步替换为 Go 后端地址。

## Mock 数据
- 文件：`web/mock-data.js`
- 已新增 `/todos` 相关 mock 以便页面在无后端时可用。
