# Catwalk CN 项目文档

## 概述

Catwalk CN 是一个 AI 推理提供商和模型的数据库服务，为 Crush 项目提供模型配置信息。本项目是对 [Charmbracelet/Catwalk](https://github.com/charmbracelet/catwalk) 的中文本地化衍生版本，保留了原始框架的核心功能，同时提供了中文文档和注释。

## 文档目录

本目录包含了 Catwalk CN 项目的完整技术文档，涵盖架构、API、开发、部署、测试、版本控制和故障排查等方面。

```
项目文档/
├── architecture/           # 架构文档
│   ├── 01-项目架构概述.md
│   ├── 02-技术栈说明.md
│   └── 03-模块划分.md
├── api/                    # API 文档
│   ├── 01-API接口文档.md
│   └── 02-数据流程文档.md
├── development/             # 开发文档
│   └── 01-开发规范.md
├── deployment/             # 部署文档
│   └── 01-构建部署流程.md
├── testing/                # 测试文档
│   └── 01-测试策略.md
├── version-control/        # 版本控制文档
│   └── 01-版本控制策略.md
└── troubleshooting/         # 故障排查文档
    └── 01-常见问题解决方案.md
```

## 文档说明

### 架构文档

#### [01-项目架构概述.md](architecture/01-项目架构概述.md)

详细描述系统整体架构设计、核心组件关系及技术架构图。

**主要内容**:
- 系统架构设计
- 核心组件说明
- 数据模型定义
- 请求流程图
- 缓存策略
- 监控指标
- 部署架构
- 扩展性设计

#### [02-技术栈说明.md](architecture/02-技术栈说明.md)

列出前端、后端、数据库、中间件等所有技术组件及其版本信息。

**主要内容**:
- 编程语言（Go 1.25.5）
- 核心依赖库
- 平台支持
- 开发工具
- HTTP 服务器
- 监控（Prometheus）
- 缓存（ETag）
- 配置管理
- 部署环境
- 安全性
- 性能优化
- 兼容性

#### [03-模块划分.md](architecture/03-模块划分.md)

明确系统功能模块划分、模块职责及模块间交互关系。

**主要内容**:
- HTTP 服务模块
- 路由处理模块
- 提供商管理模块
- 配置加载模块
- 数据模型模块
- 监控模块
- 缓存模块
- 工具模块
- 模块依赖关系
- 模块交互

### API 文档

#### [01-API接口文档.md](api/01-API接口文档.md)

提供完整的 API 接口文档，包含接口路径、请求方法、参数说明、返回格式及错误码定义。

**主要内容**:
- API 端点列表
- 请求参数说明
- 响应格式说明
- 数据模型定义
- 错误码定义
- 使用示例（cURL、Go、JavaScript、Python）
- 性能优化建议

#### [02-数据流程文档.md](api/02-数据流程文档.md)

绘制关键业务流程的数据流转图，说明数据在各模块间的传递过程。

**主要内容**:
- 系统启动流程
- 获取提供商列表流程
- 健康检查流程
- Prometheus 指标获取流程
- 数据结构流转
- 缓存流程
- 监控流程
- 错误处理流程

### 开发文档

#### [01-开发规范.md](development/01-开发规范.md)

制定编码规范、命名规范、代码审查标准及文档编写规范。

**主要内容**:
- Go 代码规范
- 命名规范
- 代码审查标准
- 文档编写规范
- Git 提交规范
- 测试规范
- 性能规范
- 安全规范
- 版本管理

### 部署文档

#### [01-构建部署流程.md](deployment/01-构建部署流程.md)

提供详细的环境配置说明、构建步骤、部署流程及环境变量配置。

**主要内容**:
- 环境要求
- 构建流程
- 跨平台构建
- Docker 构建
- 部署流程（本地、服务器、Docker、Kubernetes）
- 环境变量配置
- 反向代理配置
- 监控配置
- 健康检查
- 日志管理
- 性能优化
- 备份与恢复
- 更新流程
- 故障排查

### 测试文档

#### [01-测试策略.md](testing/01-测试策略.md)

明确单元测试、集成测试、系统测试的实施方法及测试工具使用规范。

**主要内容**:
- 测试原则
- 单元测试
- 集成测试
- 系统测试
- 测试工具
- 测试覆盖率
- 性能测试
- 测试最佳实践
- 持续集成
- 测试检查清单

### 版本控制文档

#### [01-版本控制策略.md](version-control/01-版本控制策略.md)

制定分支管理策略、代码合并流程及版本号命名规则。

**主要内容**:
- 版本控制工具
- 分支管理策略
- 分支工作流
- 版本号命名规则（语义化版本）
- 提交信息规范（约定式提交）
- 代码合并流程
- 分支保护规则
- 标签管理
- 发布流程
- 版本回退

### 故障排查文档

#### [01-常见问题解决方案.md](troubleshooting/01-常见问题解决方案.md)

整理开发、测试、部署过程中常见问题的诊断方法和解决策略。

**主要内容**:
- 开发问题（版本兼容、依赖下载、编译错误、格式化问题）
- 运行时问题（端口占用、配置加载、内存不足）
- API 问题（404、405、ETag 缓存）
- 测试问题（测试失败、覆盖率低、竞态条件）
- 部署问题（服务启动、Docker、Kubernetes）
- 性能问题（响应时间、内存泄漏、CPU 使用）
- 监控问题（Prometheus、Grafana）
- 依赖问题（依赖冲突、安全漏洞）

## 快速开始

### 1. 了解项目架构

首先阅读 [项目架构概述](architecture/01-项目架构概述.md) 了解系统整体设计。

### 2. 查看技术栈

阅读 [技术栈说明](architecture/02-技术栈说明.md) 了解项目使用的技术组件。

### 3. 了解模块划分

阅读 [模块划分](architecture/03-模块划分.md) 了解系统功能模块。

### 4. 学习 API 使用

阅读 [API 接口文档](api/01-API接口文档.md) 学习如何使用 API。

### 5. 开始开发

阅读 [开发规范](development/01-开发规范.md) 了解开发规范和最佳实践。

### 6. 部署应用

阅读 [构建部署流程](deployment/01-构建部署流程.md) 了解如何部署应用。

### 7. 编写测试

阅读 [测试策略](testing/01-测试策略.md) 了解如何编写测试。

### 8. 版本管理

阅读 [版本控制策略](version-control/01-版本控制策略.md) 了解分支管理和发布流程。

### 9. 解决问题

阅读 [常见问题解决方案](troubleshooting/01-常见问题解决方案.md) 解决遇到的问题。

## 项目信息

### 项目地址

- GitHub: https://github.com/purpose168/catwalk-cn
- 原始项目: https://github.com/charmbracelet/catwalk

### 许可证

MIT License

### 联系方式

- Email: purpose168@outlook.com

## 贡献指南

欢迎贡献代码、报告问题或提出改进建议。请遵循以下步骤：

1. Fork 项目
2. 创建功能分支（`git checkout -b feature/AmazingFeature`）
3. 提交更改（`git commit -m 'feat: 添加某个功能'`）
4. 推送到分支（`git push origin feature/AmazingFeature`）
5. 创建 Pull Request

## 更新日志

### v0.1.0 (2024-01-01)

- 初始版本发布
- 支持 21 个 AI 提供商
- 提供 RESTful API
- 集成 Prometheus 监控
- 支持 ETag 缓存

## 支持的 AI 提供商

Catwalk CN 支持以下 AI 提供商：

- OpenAI
- Anthropic
- Synthetic
- Gemini
- Azure
- Bedrock
- VertexAI
- xAI
- zAI
- OpenRouter
- GROQ
- Cerebras
- Venice
- Chutes
- DeepSeek
- HuggingFace
- AIHubMix
- Kimi Coding
- Copilot
- Vercel
- MiniMax

## 相关资源

- [Charmbracelet Catwalk](https://github.com/charmbracelet/catwalk)
- [Go 官方文档](https://golang.org/doc/)
- [Prometheus 文档](https://prometheus.io/docs/)
- [Docker 文档](https://docs.docker.com/)
- [Kubernetes 文档](https://kubernetes.io/docs/)

## 反馈与支持

如果您有任何问题或建议，请通过以下方式联系我们：

- 提交 GitHub Issue
- 发送邮件至 purpose168@outlook.com

---

**注意**: 本文档会随着项目的发展而不断更新。请定期查看最新版本。
