# 文档重组计划 - Diátaxis 框架

## 概述

本文档描述了将现有文档重组为 [Diátaxis](https://diataxis.fr/) 框架的详细计划。Diátaxis 将文档分为四类：

- **Tutorials (教程)**: 学习导向，帮助新手从 0 到 1
- **How-to Guides (操作指南)**: 目标导向，解决具体问题
- **Explanation (深度解析)**: 理解导向,讲解原理和设计
- **Reference (参考手册)**: 信息导向，提供准确的技术细节

## 目标文档结构

```
docs/
├── tutorials/                          # 教程：学习导向
│   ├── 01_getting_started.md          # 快速开始：30分钟跑起来
│   ├── 02_first_api_with_auth.md      # 第一个受保护的API
│   └── 03_extend_rbac_role.md         # 扩展RBAC角色实战
│
├── how-to/                             # 操作指南：问题导向
│   ├── configure_jwt.md               # 配置JWT参数
│   ├── manage_rbac_policies.md        # 管理RBAC策略
│   ├── protect_new_endpoint.md        # 保护新接口
│   ├── run_tests_and_ci.md            # 运行测试和CI
│   └── deploy_to_production.md        # 部署到生产环境
│
├── explanation/                        # 深度解析：理解导向
│   ├── architecture_overview.md       # 架构总览
│   ├── auth_design.md                 # JWT认证设计
│   ├── rbac_design.md                 # RBAC权限设计
│   ├── security_considerations.md     # 安全设计考量
│   ├── java_vs_go.md                  # Java vs Go深度对比 ⭐新增
│   └── design_decisions.md            # 关键设计决策
│
├── reference/                          # 参考手册：信息导向
│   ├── api.md                         # API接口参考
│   ├── auth_reference.md              # JWT字段和规范
│   ├── rbac_reference.md              # 角色权限枚举
│   ├── config.md                      # 配置项参考
│   └── error_codes.md                 # 错误码规范
│
└── meta/                               # 项目元信息
    ├── changelog.md                   # 版本变更日志
    ├── contributing.md                # 贡献指南
    └── roadmap.md                     # 开发路线图
```

## 现有文档映射方案

### 1. README.md (保留在根目录)

**当前状态**: 包含项目概述、特性列表、快速开始、API示例

**重组方案**:
- 精简为项目总入口
- 保留：项目概述、核心特性、技术栈
- 移除：详细的API示例、配置说明
- 新增：指向 Diátaxis 四类文档的导航链接

**交叉引用**:
```markdown
## 📚 文档导航

- **[教程 (Tutorials)](docs/tutorials/)** - 从零开始学习
- **[操作指南 (How-to)](docs/how-to/)** - 解决具体问题
- **[深度解析 (Explanation)](docs/explanation/)** - 理解设计原理
- **[参考手册 (Reference)](docs/reference/)** - 查阅技术细节
```

### 2. docs/AUTH.md

**当前状态**: 包含JWT认证的设计、使用、配置、安全考虑

**拆分方案**:

| 原内容章节 | 目标位置 | 文档类型 |
|-----------|---------|---------|
| JWT认证概述、设计权衡 | `explanation/auth_design.md` | Explanation |
| 如何配置JWT密钥、过期时间 | `how-to/configure_jwt.md` | How-to |
| JWT字段规范、签名算法 | `reference/auth_reference.md` | Reference |
| 快速开始示例 | `tutorials/02_first_api_with_auth.md` | Tutorial |

### 3. docs/RBAC.md

**当前状态**: 包含RBAC模型、权限矩阵、使用示例、故障排查

**拆分方案**:

| 原内容章节 | 目标位置 | 文档类型 |
|-----------|---------|---------|
| RBAC模型设计、架构说明 | `explanation/rbac_design.md` | Explanation |
| 如何添加新角色、修改权限 | `how-to/manage_rbac_policies.md` | How-to |
| 扩展RBAC实战教程 | `tutorials/03_extend_rbac_role.md` | Tutorial |
| 角色枚举、权限矩阵 | `reference/rbac_reference.md` | Reference |
| 故障排查 | `how-to/troubleshoot_rbac.md` | How-to |

### 4. docs/API.md

**当前状态**: API接口列表和请求/响应示例

**重组方案**:
- 直接迁移到 `reference/api.md`
- 保持纯参考手册风格：端点、方法、参数、响应格式
- 移除使用示例到 How-to 或 Tutorial

### 5. docs/CHANGELOG.md

**当前状态**: 版本变更记录

**重组方案**:
- 迁移到 `docs/meta/changelog.md`
- 或保留在根目录，README 中添加链接

## 新增文档规划

### ⭐ explanation/java_vs_go.md (优先级最高)

**目标读者**: 有 Java 背景的后端工程师

**章节结构**:

1. **开篇摘要 (TL;DR)**
   - 核心差异对比表格
   - 30秒快速理解

2. **语言定位与生态**
   - Java: 面向对象、企业级、JVM生态
   - Go: 简洁、并发友好、云原生

3. **运行时与部署**
   - JVM vs Go runtime
   - 启动时间、内存占用
   - 部署方式对比

4. **并发模型** ⭐核心章节
   - Thread vs Goroutine
   - 并发原语对比
   - 代码示例并排展示

5. **错误处理**
   - try-catch vs if err != nil
   - 错误处理哲学差异

6. **面向对象**
   - 类与继承 vs 结构体与组合
   - 接口实现方式对比

7. **Web框架与生态**
   - Spring Boot vs Gin/net/http
   - 依赖注入、中间件模式

8. **安全与认证实践** (结合本项目)
   - Spring Security vs Go中间件
   - JWT + RBAC实现对比

9. **迁移策略**
   - 从Java迁移到Go的路径
   - 常见陷阱和最佳实践

**排版设计**:
- 使用 "For Java Devs" 提示框
- 代码示例并排对比
- 关键概念用表格总结

### tutorials/01_getting_started.md

**目标**: 30分钟内让新手跑起来项目

**内容**:
1. 环境准备 (Go安装、IDE配置)
2. 克隆项目
3. 配置数据库
4. 启动服务
5. 测试第一个API请求
6. 下一步学习路径

### tutorials/02_first_api_with_auth.md

**目标**: 创建第一个受JWT+RBAC保护的API

**内容**:
1. 注册用户
2. 登录获取Token
3. 使用Token访问受保护接口
4. 理解认证流程
5. 常见错误排查

### how-to/configure_jwt.md

**目标**: 配置JWT参数

**内容**:
- 修改JWT密钥
- 调整Token过期时间
- 配置刷新Token策略
- 环境变量配置

### reference/error_codes.md (新增)

**目标**: 统一的错误码规范

**内容**:
- HTTP状态码映射
- 业务错误码定义
- 错误响应格式
- 错误处理最佳实践

## 交叉引用策略

### 1. 标准化提示框

使用 Markdown 扩展语法或统一的格式:

```markdown
> **[深入理解]** 想了解JWT的签名算法选择，请阅读 [JWT认证设计](../explanation/auth_design.md#签名算法)

> **[查阅参考]** 完整的API列表请参考 [API参考手册](../reference/api.md)

> **[前置知识]** 本指南假设您已完成 [快速开始教程](../tutorials/01_getting_started.md)

> **[动手实践]** 查看 [如何保护新接口](../how-to/protect_new_endpoint.md) 了解实际应用
```

### 2. 文档间链接规范

- 使用相对路径
- 包含锚点定位到具体章节
- 在链接文本中说明目标文档类型

### 3. 导航设计

每个文档开头包含面包屑导航:

```markdown
[首页](../../README.md) > [教程](../tutorials/) > 快速开始
```

每个文档结尾包含"下一步"推荐:

```markdown
## 下一步

- 📖 [创建第一个受保护的API](02_first_api_with_auth.md)
- 🔍 [深入理解JWT认证设计](../explanation/auth_design.md)
- 📚 [API参考手册](../reference/api.md)
```

## 视觉设计规范

### 代码示例

```markdown
**文件**: `internal/todo/auth.go`

\`\`\`go
// (1) 创建JWT管理器
jwtManager := NewJWTManager("secret", 24*time.Hour)

// (2) 生成Token
token, err := jwtManager.Generate(userID)
if err != nil {
    return err
}
\`\`\`

**说明**:
1. 使用密钥和过期时间初始化管理器
2. 为指定用户生成Token
```

### 对比表格

```markdown
| 特性 | Java | Go |
|------|------|-----|
| 并发模型 | Thread + ThreadPool | Goroutine + Channel |
| 错误处理 | try-catch-finally | if err != nil |
| 部署方式 | JAR/WAR + JVM | 单一二进制文件 |
```

### 图表

使用 Mermaid.js:

```markdown
\`\`\`mermaid
sequenceDiagram
    Client->>+Server: POST /login
    Server->>+UserStore: FindByEmail
    UserStore-->>-Server: User
    Server->>Server: CheckPassword
    Server->>+JWTManager: Generate(userID)
    JWTManager-->>-Server: token
    Server-->>-Client: {token, expires_in}
\`\`\`

**图1**: JWT登录流程时序图
```

## 实施步骤

### 阶段1: 规划与设计 ✅ (当前阶段)

- [x] 确定Diátaxis目录结构
- [x] 制定现有文档映射方案
- [x] 设计交叉引用策略
- [x] 规划视觉设计规范

### 阶段2: 创建新文档

1. **优先**: 创建 `explanation/java_vs_go.md`
2. 创建 `tutorials/01_getting_started.md`
3. 创建 `tutorials/02_first_api_with_auth.md`
4. 创建 `reference/error_codes.md`

### 阶段3: 拆分现有文档

1. 拆分 `docs/AUTH.md` 到三个目标位置
2. 拆分 `docs/RBAC.md` 到四个目标位置
3. 迁移 `docs/API.md` 到 `reference/api.md`
4. 迁移 `docs/CHANGELOG.md` 到 `meta/changelog.md`

### 阶段4: 更新README和导航

1. 精简 README.md
2. 添加Diátaxis导航链接
3. 在每个文档中添加面包屑和"下一步"

### 阶段5: 审查与优化

1. 检查所有交叉引用链接
2. 统一代码示例格式
3. 添加图表和表格
4. 与codex进行最终审查

## 成功标准

- [ ] 所有文档归类到Diátaxis四个象限
- [ ] 交叉引用链接完整且准确
- [ ] 代码示例可运行且有注释
- [ ] 图表清晰且统一风格
- [ ] Java开发者能快速找到对比指南
- [ ] 新手能通过Tutorial在30分钟内跑起来

## 参考资源

- [Diátaxis官方文档](https://diataxis.fr/)
- [Write the Docs](https://www.writethedocs.org/)
- [Google Developer Documentation Style Guide](https://developers.google.com/style)

---

**创建时间**: 2025-12-05
**状态**: 规划完成，待执行
