# 001 - JWT + RBAC + Refresh Token 设计决策

## 背景
- 项目已支持 JWT 登录与基础 RBAC，中短期需求是补齐刷新令牌、登录失败限制及前端联动。
- 场景：教学与演示为主，但仍需体现安全基线与可扩展性。

## 决策
1. **令牌模型**：Access Token 使用 HS256 JWT，TTL 默认 24h；Refresh Token 使用 32 字节随机值存储在内存表，TTL 7d，单用户仅保留一枚有效 refresh（旋转时作废旧 token）。  
2. **RBAC**：沿用 Role → Permission 映射（admin/user/guest），在授权中间件检查资源/动作并验证归属。  
3. **登录失败限制**：同一邮箱 15 分钟窗口内最多 5 次失败，超过锁定 10 分钟，返回 429 + Retry-After。成功登录清空计数。  
4. **刷新接口**：`POST /refresh` 接收 refresh_token，校验、旋转并返回新 access+refresh。  
5. **前端联动**：新增 `web/auth-helper.js` 提供 `authFetch` 自动刷新，`web/todo-login.html` 作为演示页。

## 备选方案
- 将 refresh 也做成 JWT（带 jti + 队列），但教学场景下内存表更直观且无外部存储依赖。
- 使用 IP 级限流替代邮箱维度，但邮箱维度更贴合账户暴力破解防护。

## 影响
- 内存表实现简洁但非分布式；生产需要持久化（Redis）并支持多节点同步。
- 单用户单 refresh 的设计简化撤销逻辑，但无法多终端并存；教学足够。

## TODO
- 将 refresh/session 存入 Redis，支持多实例。
- 补充前后端 E2E 测试覆盖 refresh 流程。
- 在网关侧复用同一 RBAC/鉴权逻辑，避免重复实现。
