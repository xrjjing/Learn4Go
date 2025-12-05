# 故障排查指南

## 登录相关
- **症状：401 authorization required**  
  - 检查是否携带 `Authorization: Bearer <token>`  
  - token 过期：调用 `POST /refresh` 获取新 access。
- **症状：429 account temporarily locked**  
  - 15 分钟内失败超过 5 次会锁定 10 分钟；查看响应头 `Retry-After`。  
  - 解决：等待或重置密码，成功登录会清空计数。

## 刷新令牌
- **症状：401 refresh token expired/invalid**  
  - refresh 仅保留一枚有效；换端登录会使旧 refresh 失效。  
  - 超过 7 天需重新登录。

## RBAC
- **症状：403 insufficient permissions / you don't own this resource**  
  - 角色权限不够或资源归属不匹配。  
  - 管理员可操作全部；普通用户只能操作自己的 TODO。

## 服务与依赖
- **后端未响应 /healthz**：确认 `cmd/todoapi` 进程在 8080；本地可运行 `./start-local.sh`。  
- **Redis/MinIO 相关错误**：启用对应模块前确保 Docker 依赖已启动（见 `deployments/docker-compose.yml`）。

## 调试建议
- 开启调试日志：查看 `cmd/todoapi` 控制台输出。  
- 使用 `curl -v` 检查头部与返回码。  
- 若定位性能问题，开启 `net/http/pprof`（可在 main 中临时引入）。
