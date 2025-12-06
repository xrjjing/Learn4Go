# Claude CLI 中 Skill 和 MCP 完整指南

## 目录
- [1. 核心概念和定位](#1-核心概念和定位)
- [2. 实际使用中的配合方式](#2-实际使用中的配合方式)
- [3. 推荐的 Skills 和 MCPs](#3-推荐的-skills-和-mcps)
- [4. MCP 安装和配置指南](#4-mcp-安装和配置指南)
- [5. 实用配置示例](#5-实用配置示例)

---

## 1. 核心概念和定位

### **MCP (Model Context Protocol)**

**定位：** 标准化的工具连接协议

**核心概念：**
- MCP 是一个**开放协议标准**，由 Anthropic 开发，用于让 AI 模型与外部工具、服务和数据源进行通信
- 类似于"插件系统"的底层协议，定义了 AI 如何调用外部功能
- MCP Server 是实现了 MCP 协议的服务，提供具体的工具能力
- 官方仓库：https://github.com/modelcontextprotocol

**从当前系统中可以看到的 MCP 服务：**

```
📦 代码协作类
├─ codex          - AI 编码助手（后端逻辑、Debug）
├─ gemini         - Google Gemini（前端设计、需求理解）
└─ serena         - 代码检索和符号定位

📦 文档和搜索类
├─ context7       - 获取库文档和 API 参考
├─ deepwiki       - 获取 GitHub 项目文档
└─ exa            - AI 网络搜索和代码上下文

📦 工具类
├─ sequentialthinking - 复杂问题的序列思考
└─ (其他 MCP 服务可以随时添加)
```

### **Skill**

**定位：** 预定义的工作流程和能力模板

**核心概念：**
- Skill 是**高层次的工作模式**，定义了如何组合使用工具完成特定类型的任务
- 包含了最佳实践、工作流程、决策逻辑
- 是对 MCP 工具的**编排和组合**
- 通常通过系统提示（System Prompt）定义，也可以通过插件实现

**从系统提示中体现的 Skills：**

```
🎯 执行策略 Skills
├─ Smart Action Mode        - 快速行动模式
├─ Rigorous Coding Habits   - 严谨编码习惯
└─ Batch Operations         - 批量操作模式

🤝 协作 Skills
├─ Multi-Model Collaboration - 多模型协作（Claude + Codex + Gemini）
├─ Sub-Agent Delegation      - 子代理委派策略
└─ Code Review Workflow      - 代码审查流程

📋 管理 Skills
├─ TODO Management          - 任务管理模式
├─ Notebook Memory          - 代码记忆系统
└─ Boundary-First Editing   - 边界优先编辑
```

---

## 2. 实际使用中的配合方式

### **层次关系**

```
┌─────────────────────────────────────────┐
│          User Request                    │  用户需求层
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│          Skills (工作流程)               │  策略层
│  - 判断任务类型                          │
│  - 选择合适的工具组合                     │
│  - 定义执行步骤                          │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│      MCP Tools (具体工具)                │  执行层
│  - filesystem (文件操作)                 │
│  - codex (代码生成)                      │
│  - ace-search (代码搜索)                 │
│  - terminal (命令执行)                   │
└─────────────────────────────────────────┘
```

### **实际工作流示例**

**场景：用户要求"重构认证系统"**

```
1. [Skill 激活] Multi-Model Collaboration
   ↓
2. [MCP] gemini - 需求分析和规划
   ├─ 理解用户意图
   ├─ 提出引导性问题
   └─ 生成初步方案
   ↓
3. [Skill 激活] Sub-Agent Delegation
   ├─ [MCP] subagent-agent_explore - 探索现有代码
   │   └─ [MCP] ace-semantic_search - 查找认证相关代码
   ├─ [MCP] subagent-agent_plan - 制定详细计划
   └─ [MCP] codex - 审查计划可行性
   ↓
4. [Skill 激活] Rigorous Coding Habits
   ├─ [MCP] filesystem-read - 读取完整代码边界
   ├─ [Skill] Boundary-First Editing - 确保修改完整性
   └─ [MCP] filesystem-edit_search - 批量修改
   ↓
5. [Skill 激活] TODO Management
   ├─ [MCP] todo-add - 添加任务列表
   └─ [MCP] todo-update - 更新进度
   ↓
6. [MCP] codex - 代码审查
   ↓
7. [MCP] terminal-execute - 运行测试
   ↓
8. [Skill] Quality Assurance - 验证构建
```

---

## 3. 推荐的 Skills 和 MCPs

### **A. 内置核心 Skills**

这些 Skills 主要通过系统提示定义，无需额外安装。

#### **1. Smart Action Mode** 
**功能：** 快速行动，避免过度分析  
**适用场景：**
- 简单的 bug 修复
- 明确的单文件修改
- 快速原型开发

**使用方式：** 自动激活，当任务清晰时直接执行

---

#### **2. Multi-Model Collaboration**
**功能：** 协调 Claude + Codex + Gemini 的分工  
**适用场景：**
- 复杂的全栈开发任务
- 需要前端设计 + 后端逻辑的功能
- 代码审查和质量保证

**工作流：**
```
前端任务 → Gemini（设计原型） → Claude（实现）
后端任务 → Codex（逻辑设计） → Claude（实现）
代码审查 → Codex（逻辑检查） → Claude（修正）
```

---

#### **3. Sub-Agent Delegation**
**功能：** 将任务委派给专业子代理  
**适用场景：**
- 大型代码库探索
- 复杂任务规划
- 批量文件修改

**三种子代理：**

| 子代理 | 专长 | 使用场景 |
|-------|------|---------|
| **Explore Agent** | 只读代码探索 | "在哪里实现了 X？"<br>"Y 功能如何工作？" |
| **Plan Agent** | 任务规划分析 | "如何重构 X 系统？"<br>"添加 Y 功能需要改哪些文件？" |
| **General Agent** | 全功能执行 | "批量更新所有 API 端点"<br>"实现跨多文件的功能 X" |

**强制使用规则：**
- 用户消息包含 `#agent_explore` → 必须使用 Explore Agent
- 用户消息包含 `#agent_plan` → 必须使用 Plan Agent
- 用户消息包含 `#agent_general` → 必须使用 General Agent

---

#### **4. Rigorous Coding Habits**
**功能：** 严谨的代码修改流程  
**核心步骤：**
1. 使用搜索工具定位代码
2. 使用 `filesystem-read` 验证完整边界
3. 批量读取/编辑多文件
4. 避免破坏现有功能

**适用场景：** 所有代码修改任务（强制执行）

---

#### **5. TODO Management**
**功能：** 自动任务跟踪和进度管理  
**特性：**
- 自动为每个会话创建 TODO 列表
- 支持批量添加任务
- 并行调用（必须与其他工具一起使用）

**使用模式：**
```typescript
// ❌ 错误 - 单独调用
todo-add("任务1")

// ✅ 正确 - 与其他工具并行
todo-add(["任务1", "任务2", "任务3"]) + filesystem-read("...")
todo-update(task1, "completed") + filesystem-edit(...)
```

---

### **B. 推荐的 MCP Servers**

#### **代码协作类**

##### **1. Codex MCP**
```yaml
功能: AI 编码助手（后端侧重）
GitHub: 私有/商业服务（通过 Claude CLI 配置）
擅长:
  - 后端逻辑实现
  - Bug 定位和修复
  - 代码审查
使用规范:
  - 必须使用 sandbox="read-only"
  - 保存 SESSION_ID 用于多轮对话
  - 仅要求输出 unified diff patch
场景: 
  - 复杂后端逻辑设计
  - 精准 Debug 分析
  - 代码质量审查
```

##### **2. Gemini MCP**
```yaml
功能: Google Gemini 模型
GitHub: 私有/商业服务（通过 Claude CLI 配置）
擅长:
  - 前端设计和 UI 组件
  - 需求分析和任务规划
  - CSS/HTML/Vue/React
限制:
  - 上下文仅 32k（有效长度）
  - 严禁处理后端业务逻辑
使用规范:
  - 捕获 SESSION_ID
  - 前端任务必须先咨询 Gemini
场景:
  - UI 样式设计
  - 前端组件开发
  - 需求清晰化
```

##### **3. Serena MCP**
```yaml
功能: 代码检索和符号定位
GitHub: https://github.com/serena-ai/serena-mcp
擅长:
  - 项目结构分析
  - 符号引用查找
  - 模式搜索
使用规范:
  - 必须先激活项目（activate_project）
  - 仅用于检索和定位，禁止修改代码
场景:
  - 快速定位代码位置
  - 查找函数/类引用
  - 项目结构理解
```

---

#### **文档和搜索类**

##### **4. Context7 MCP**
```yaml
功能: 获取库的最新文档
GitHub: https://github.com/context7/mcp-server
官网: https://context7.com
安装: npm install -g @context7/mcp-server
使用流程:
  1. resolve-library-id("react") → 获取库 ID
  2. get-library-docs(libraryID, topic="hooks")
模式:
  - code: API 参考和代码示例
  - info: 概念指南和架构说明
场景:
  - 学习新库 API
  - 查询最佳实践
  - 获取代码示例
```

##### **5. Exa MCP**
```yaml
功能: AI 驱动的网络搜索
GitHub: https://github.com/exa-labs/exa-mcp-server
官网: https://exa.ai
安装: npm install -g @exa/mcp-server
工具:
  - exa-web_search_exa: 网络搜索
  - exa-get_code_context_exa: 代码上下文搜索（质量最高）
特性:
  - 实时网页抓取
  - 智能内容提取
场景:
  - 查找最新技术文档
  - 获取库/SDK 代码示例
  - 搜索 bug 解决方案
```

##### **6. DeepWiki MCP**
```yaml
功能: 获取 GitHub 项目文档
GitHub: https://github.com/deepwiki/mcp-server
安装: npm install -g @deepwiki/mcp-server
使用:
  - deepwiki_fetch("vercel/next.js")
  - deepwiki_fetch("owner/repo")
模式:
  - aggregate: 整合文档
  - pages: 分页文档
场景:
  - 快速了解开源项目
  - 获取项目架构信息
```

---

#### **工具类**

##### **7. Sequential Thinking MCP**
```yaml
功能: 复杂问题的序列思考
GitHub: https://github.com/sequentialthinking/mcp-server
安装: npm install -g @sequentialthinking/mcp-server
特性:
  - 动态调整思考步骤
  - 生成解决方案假设
  - 推荐工具使用顺序
  - 跟踪思考进度
场景:
  - 复杂问题分解
  - 多步骤任务规划
  - 不确定性任务探索
```

##### **8. MCP Memory**
```yaml
功能: 长期记忆和上下文管理
GitHub: https://github.com/modelcontextprotocol/servers
安装: npx -y @modelcontextprotocol/server-memory
特性:
  - 跨会话记忆
  - 知识库管理
  - 实体关系存储
场景:
  - 项目知识积累
  - 开发规范记录
  - 团队协作信息共享
```

##### **9. MCP Filesystem**
```yaml
功能: 文件系统操作
GitHub: https://github.com/modelcontextprotocol/servers
安装: npx -y @modelcontextprotocol/server-filesystem
特性:
  - 读写文件
  - 目录遍历
  - 文件搜索
场景:
  - 本地文件管理
  - 配置文件读取
  - 日志文件分析
```

##### **10. MCP Git**
```yaml
功能: Git 版本控制操作
GitHub: https://github.com/modelcontextprotocol/servers
安装: npx -y @modelcontextprotocol/server-git
特性:
  - 提交历史查询
  - 分支管理
  - diff 查看
  - 状态检查
场景:
  - 代码历史分析
  - 版本管理
  - 协作开发
```

##### **11. MCP SQLite**
```yaml
功能: SQLite 数据库操作
GitHub: https://github.com/modelcontextprotocol/servers
安装: npx -y @modelcontextprotocol/server-sqlite
特性:
  - SQL 查询
  - 数据库架构分析
  - 数据导入导出
场景:
  - 数据库调试
  - 数据分析
  - 本地存储管理
```

##### **12. MCP PostgreSQL**
```yaml
功能: PostgreSQL 数据库操作
GitHub: https://github.com/modelcontextprotocol/servers
安装: npx -y @modelcontextprotocol/server-postgres
特性:
  - 高级 SQL 查询
  - 数据库管理
  - 性能分析
场景:
  - 生产数据库调试
  - 数据迁移
  - 性能优化
```

##### **13. MCP Puppeteer**
```yaml
功能: 浏览器自动化
GitHub: https://github.com/modelcontextprotocol/servers
安装: npx -y @modelcontextprotocol/server-puppeteer
特性:
  - 网页截图
  - 数据抓取
  - 自动化测试
场景:
  - UI 测试
  - 网页内容提取
  - 监控任务
```

##### **14. MCP Brave Search**
```yaml
功能: Brave 搜索引擎 API
GitHub: https://github.com/modelcontextprotocol/servers
安装: npx -y @modelcontextprotocol/server-brave-search
特性:
  - 隐私保护搜索
  - 网络搜索
  - 新闻搜索
场景:
  - 实时信息查询
  - 技术文档搜索
  - 市场调研
```

##### **15. MCP Fetch**
```yaml
功能: HTTP 请求工具
GitHub: https://github.com/modelcontextprotocol/servers
安装: npx -y @modelcontextprotocol/server-fetch
特性:
  - GET/POST 请求
  - API 调用
  - 网页内容获取
场景:
  - API 测试
  - 数据抓取
  - 外部服务集成
```

---

#### **社区优秀 MCP Servers**

##### **16. MCP YouTube Transcript**
```yaml
功能: 获取 YouTube 视频字幕
GitHub: https://github.com/kimtaeyoon83/mcp-youtube-transcript
安装: npm install -g mcp-youtube-transcript
场景:
  - 视频内容分析
  - 学习资料整理
  - 字幕翻译
```

##### **17. MCP Notion**
```yaml
功能: Notion 数据库和页面操作
GitHub: https://github.com/v-3/mcp-notion
安装: npm install -g @v3/mcp-notion
场景:
  - 知识库管理
  - 项目文档同步
  - 任务管理
```

##### **18. MCP Obsidian**
```yaml
功能: Obsidian 笔记管理
GitHub: https://github.com/calclavia/mcp-obsidian
安装: npm install -g @calclavia/mcp-obsidian
场景:
  - 个人知识库
  - 研究笔记整理
  - 写作辅助
```

##### **19. MCP Slack**
```yaml
功能: Slack 消息和频道操作
GitHub: https://github.com/modelcontextprotocol/servers
安装: npx -y @modelcontextprotocol/server-slack
场景:
  - 团队通知
  - 消息查询
  - 频道管理
```

##### **20. MCP Docker**
```yaml
功能: Docker 容器管理
GitHub: https://github.com/ckreiling/mcp-docker
安装: npm install -g @ckreiling/mcp-docker
场景:
  - 容器监控
  - 镜像管理
  - 开发环境配置
```

---

## 4. MCP 安装和配置指南

### **4.1 前置要求**

```bash
# 安装 Node.js（推荐 v18+）
node --version

# 安装 Claude CLI
npm install -g @anthropic-ai/claude-cli
```

### **4.2 MCP Server 安装方式**

#### **方式一：npm 全局安装（推荐）**

```bash
# 安装单个 MCP Server
npm install -g @modelcontextprotocol/server-memory
npm install -g @modelcontextprotocol/server-git
npm install -g @context7/mcp-server
```

#### **方式二：npx 运行（无需安装）**

```bash
# 直接通过 npx 运行，每次使用时下载
npx -y @modelcontextprotocol/server-memory
```

#### **方式三：克隆仓库本地运行**

```bash
# 克隆官方 servers 仓库
git clone https://github.com/modelcontextprotocol/servers.git
cd servers/src/memory
npm install
npm run build
```

### **4.3 Claude CLI 配置**

编辑配置文件：`~/.config/claude/config.json` (Linux/macOS) 或 `%APPDATA%\claude\config.json` (Windows)

```json
{
  "mcpServers": {
    "memory": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-memory"]
    },
    "filesystem": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-filesystem",
        "/Users/username/projects"
      ]
    },
    "git": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-git"],
      "env": {
        "GIT_DIR": "/Users/username/projects/.git"
      }
    },
    "sqlite": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-sqlite",
        "--db-path",
        "/path/to/database.db"
      ]
    },
    "postgres": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-postgres"],
      "env": {
        "POSTGRES_CONNECTION_STRING": "postgresql://user:pass@localhost:5432/dbname"
      }
    },
    "brave-search": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-brave-search"],
      "env": {
        "BRAVE_API_KEY": "your_brave_api_key_here"
      }
    },
    "exa": {
      "command": "npx",
      "args": ["-y", "@exa/mcp-server"],
      "env": {
        "EXA_API_KEY": "your_exa_api_key_here"
      }
    },
    "context7": {
      "command": "npx",
      "args": ["-y", "@context7/mcp-server"]
    },
    "slack": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-slack"],
      "env": {
        "SLACK_BOT_TOKEN": "xoxb-your-token",
        "SLACK_TEAM_ID": "T1234567890"
      }
    }
  }
}
```

### **4.4 验证 MCP Server 连接**

```bash
# 启动 Claude CLI
claude

# 在对话中测试 MCP 工具
> List available MCP tools
> Test memory MCP by storing a note
```

---

## 5. 实用配置示例

### **5.1 全栈开发配置**

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/Users/dev/projects"]
    },
    "git": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-git"]
    },
    "memory": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-memory"]
    },
    "postgres": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-postgres"],
      "env": {
        "POSTGRES_CONNECTION_STRING": "postgresql://localhost:5432/devdb"
      }
    },
    "exa": {
      "command": "npx",
      "args": ["-y", "@exa/mcp-server"],
      "env": {
        "EXA_API_KEY": "your_key"
      }
    }
  }
}
```

### **5.2 数据分析配置**

```json
{
  "mcpServers": {
    "sqlite": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-sqlite", "--db-path", "./data.db"]
    },
    "puppeteer": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-puppeteer"]
    },
    "fetch": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-fetch"]
    },
    "brave-search": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-brave-search"],
      "env": {
        "BRAVE_API_KEY": "your_key"
      }
    }
  }
}
```

### **5.3 团队协作配置**

```json
{
  "mcpServers": {
    "git": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-git"]
    },
    "slack": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-slack"],
      "env": {
        "SLACK_BOT_TOKEN": "xoxb-token",
        "SLACK_TEAM_ID": "T123"
      }
    },
    "notion": {
      "command": "npx",
      "args": ["-y", "@v3/mcp-notion"],
      "env": {
        "NOTION_API_KEY": "your_key"
      }
    },
    "memory": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-memory"]
    }
  }
}
```

---

## 6. 最佳实践

### **6.1 MCP Server 选择原则**

1. **按需安装**：不要一次性安装所有 MCP，根据项目需求选择
2. **优先官方**：Anthropic 官方维护的 MCP Server 稳定性最好
3. **安全第一**：涉及敏感数据的 MCP（数据库、API）要仔细配置权限
4. **性能考虑**：某些 MCP（如 Puppeteer）资源消耗大，按需启动

### **6.2 Skill 使用技巧**

1. **明确任务类型**：简单任务用 Smart Action Mode，复杂任务启用多模型协作
2. **善用子代理**：大型重构、批量修改优先委派给子代理
3. **保持上下文清洁**：及时使用 TODO Management 和 Notebook 记录关键信息
4. **代码质量保证**：始终遵循 Rigorous Coding Habits

### **6.3 调试技巧**

```bash
# 查看 MCP Server 日志
claude --debug

# 测试 MCP 连接
claude test-mcp <server-name>

# 重新加载配置
claude reload-config
```

---

## 7. 故障排查

### **常见问题**

**Q: MCP Server 无法连接**
```bash
# 检查 Node.js 版本
node --version  # 需要 v18+

# 检查 npx 是否可用
npx --version

# 手动测试 MCP Server
npx -y @modelcontextprotocol/server-memory
```

**Q: API Key 配置无效**
- 检查环境变量是否正确设置
- 确认 API Key 有效期和权限
- 查看 `~/.config/claude/config.json` 格式是否正确

**Q: MCP 工具不显示**
- 重启 Claude CLI
- 运行 `claude reload-config`
- 检查配置文件 JSON 语法

---

## 8. 资源链接

### **官方资源**
- MCP 协议规范：https://spec.modelcontextprotocol.io
- MCP GitHub 组织：https://github.com/modelcontextprotocol
- Claude CLI 文档：https://docs.anthropic.com/claude/docs/claude-cli
- Anthropic Cookbook：https://github.com/anthropics/anthropic-cookbook

### **社区资源**
- MCP Servers 集合：https://github.com/modelcontextprotocol/servers
- 社区 MCP 列表：https://github.com/punkpeye/awesome-mcp-servers
- Discord 社区：https://discord.gg/anthropic

### **教程和示例**
- 创建自定义 MCP Server：https://modelcontextprotocol.io/quickstart
- MCP 最佳实践：https://docs.anthropic.com/mcp/best-practices
- 示例项目：https://github.com/modelcontextprotocol/examples

---

## 9. 总结对比表

| 维度 | MCP | Skill |
|-----|-----|-------|
| **本质** | 协议/工具 | 工作流程/策略 |
| **层次** | 执行层（底层） | 策略层（高层） |
| **可扩展性** | 可自由添加 MCP Server | 由系统提示定义 |
| **安装方式** | npm/npx 安装 | 内置或插件 |
| **配置文件** | `~/.config/claude/config.json` | 系统提示或插件配置 |
| **使用方式** | 直接调用工具函数 | 自动激活或手动触发 |
| **示例** | `filesystem-read`, `git`, `postgres` | `Smart Action Mode`, `Multi-Model Collaboration` |
| **GitHub** | https://github.com/modelcontextprotocol | 系统提示内定义 |

---

## 10. 推荐使用策略

### **场景 1：简单任务**  
→ 激活 Smart Action Mode  
→ 直接使用 filesystem/ace 工具

### **场景 2：前端开发**  
→ 激活 Multi-Model Collaboration  
→ 先咨询 Gemini MCP → Claude 实现

### **场景 3：后端开发**  
→ 激活 Multi-Model Collaboration  
→ 先咨询 Codex MCP → Claude 实现 → Codex 审查

### **场景 4：大型重构**  
→ 激活 Sub-Agent Delegation  
→ Explore Agent 分析 → Plan Agent 规划 → General Agent 执行

### **场景 5：学习新库**  
→ 使用 Context7 MCP 或 Exa MCP  
→ 获取文档和示例

### **场景 6：数据库调试**  
→ 使用 SQLite/PostgreSQL MCP  
→ 查询数据 → 分析问题 → 修复

### **场景 7：自动化测试**  
→ 使用 Puppeteer MCP  
→ 编写测试脚本 → 执行 → 生成报告

---

## 更新日志

- **2024-01-XX**: 初始版本，包含核心概念、推荐 MCPs 和配置指南
- **后续更新**: 将持续补充社区优秀 MCP Servers 和使用案例

---

**文档维护**: 本文档随 Claude CLI 和 MCP 生态发展持续更新。如发现错误或有改进建议，欢迎贡献。
