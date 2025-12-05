# 架构草图（Go 版）

目标：对标 Python 仓库的「前端 + 网关 + 认证 + 业务/日志」链路，使用 Go 标准库起步，后续可引入 Gin/Fiber。前端静态页位于 `web/`。

## 拆分建议
- **gateway_service**：JWT 校验 + 统一响应包装 + 反向代理（标准库 `net/http/httputil`）
- **auth_service**：最小 JWT 签发与校验（`github.com/golang-jwt/jwt/v5`，若网络受限可先用伪实现）
- **user_order_service**：用户/订单查询与创建（内存或 SQLite）
- **log_service**：日志收集与分析（正则 + 规则匹配）
- **frontend**：`web/` 目录静态页，指向 gateway 路由

## 技术选型
- Web：标准库 `net/http`；若可联网再切 Gin/Echo 以提升 DX
- 配置：`flag` + 环境变量，预留 `spf13/viper` 入口
- 日志：`log/slog` 或 `zap`（需网络）
- 数据：示例使用内存 map，后续可接 SQLite (`modernc.org/sqlite` 可离线) 或 MySQL 驱动
- 并发：goroutine + channel + context，用中间件式取消/超时

## 路由示例（标准库版）
- `POST /auth/login` → 签发伪 JWT
- `GET /api/users` → 返回示例用户列表
- `POST /api/orders` → 创建订单
- `POST /log/analyze` → 提交日志文本，返回匹配结果

## 后续演进
1. 引入 Gin：路由分组/中间件、绑定校验、swagger 文档
2. 增加统一错误码与响应包装
3. gRPC 版本的 user_order_service + gateway 转码
4. Prometheus 指标导出与 pprof 性能分析
