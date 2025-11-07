# 贡献指南 / Contributing Guide

感谢您对 HomeGuard 项目的关注！

Thank you for your interest in contributing to HomeGuard!

[中文](#中文) | [English](#english)

---

## 中文

### 如何贡献

我们欢迎任何形式的贡献，包括但不限于：

- 🐛 报告 Bug
- 💡 提出新功能建议
- 📝 改进文档
- 🔧 提交代码修复或新功能
- 🌍 翻译文档

### 开发环境设置

1. **安装 Go**
   - 需要 Go 1.23 或更高版本
   - 下载：https://golang.org/dl/

2. **克隆仓库**
   ```bash
   git clone https://github.com/p3ddd/HomeGuard.git
   cd HomeGuard
   ```

3. **安装依赖**
   ```bash
   go mod download
   ```

4. **运行项目**
   ```bash
   go run main.go
   ```

### 开发流程

1. **Fork 项目**
   - 点击页面右上角的 "Fork" 按钮

2. **创建特性分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或者修复分支
   git checkout -b fix/bug-description
   ```

3. **开发和测试**
   ```bash
   # 运行测试
   go test ./...
   
   # 运行 linter
   golangci-lint run
   
   # 格式化代码
   go fmt ./...
   ```

4. **提交代码**
   ```bash
   git add .
   git commit -m "feat: 添加新功能描述"
   # 或
   git commit -m "fix: 修复某个问题"
   ```

5. **推送到 GitHub**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **创建 Pull Request**
   - 在 GitHub 上创建 Pull Request
   - 详细描述您的更改
   - 等待审核

### 代码规范

#### Commit 消息格式

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型 (type)：**
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响代码运行）
- `refactor`: 重构（既不是新功能也不是 Bug 修复）
- `test`: 添加或修改测试
- `chore`: 构建过程或辅助工具的变动

**示例：**
```
feat(mqtt): 添加 MQTT QoS 配置支持

- 允许用户配置 MQTT QoS 级别
- 添加命令行参数 -mqtt-qos
- 更新文档

Closes #123
```

#### 代码风格

- 使用 `go fmt` 格式化代码
- 使用 `golangci-lint` 检查代码质量
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 指南
- 为导出的函数和类型添加注释
- 保持函数简短，职责单一

#### 测试要求

- 为新功能添加单元测试
- 确保所有测试通过：`go test ./...`
- 维持或提高代码覆盖率

### 项目结构

```
HomeGuard/
├── device/           # 设备管理模块
├── listener/         # 监听器实现
├── wol/             # WOL 核心功能
├── main.go          # 程序入口
└── ...
```

### 报告 Bug

在提交 Bug 报告时，请包含：

1. **问题描述**：清晰简洁地描述问题
2. **复现步骤**：详细的复现步骤
3. **期望行为**：您期望发生什么
4. **实际行为**：实际发生了什么
5. **环境信息**：
   - 操作系统
   - Go 版本
   - HomeGuard 版本
6. **日志输出**：相关的日志或错误信息

### 功能建议

在提交功能建议时，请包含：

1. **问题或需求**：描述您想解决的问题
2. **建议的解决方案**：您的想法
3. **替代方案**：您考虑过的其他方案
4. **使用场景**：该功能的具体应用场景

### 行为准则

- 尊重所有贡献者
- 包容不同的观点和经验
- 接受建设性的批评
- 专注于对社区最有利的事情

---

## English

### How to Contribute

We welcome all forms of contributions, including but not limited to:

- 🐛 Reporting bugs
- 💡 Suggesting new features
- 📝 Improving documentation
- 🔧 Submitting bug fixes or new features
- 🌍 Translating documentation

### Development Setup

1. **Install Go**
   - Go 1.23 or higher required
   - Download: https://golang.org/dl/

2. **Clone Repository**
   ```bash
   git clone https://github.com/p3ddd/HomeGuard.git
   cd HomeGuard
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Run Project**
   ```bash
   go run main.go
   ```

### Development Workflow

1. **Fork the Project**
   - Click the "Fork" button in the top right

2. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or fix branch
   git checkout -b fix/bug-description
   ```

3. **Develop and Test**
   ```bash
   # Run tests
   go test ./...
   
   # Run linter
   golangci-lint run
   
   # Format code
   go fmt ./...
   ```

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   # or
   git commit -m "fix: fix some issue"
   ```

5. **Push to GitHub**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create Pull Request**
   - Create a Pull Request on GitHub
   - Describe your changes in detail
   - Wait for review

### Code Standards

#### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation updates
- `style`: Code formatting (no code logic change)
- `refactor`: Code refactoring
- `test`: Adding or modifying tests
- `chore`: Build process or tool changes

**Example:**
```
feat(mqtt): add MQTT QoS configuration support

- Allow users to configure MQTT QoS level
- Add command line parameter -mqtt-qos
- Update documentation

Closes #123
```

#### Code Style

- Format code with `go fmt`
- Check code quality with `golangci-lint`
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Add comments for exported functions and types
- Keep functions short and focused

#### Testing Requirements

- Add unit tests for new features
- Ensure all tests pass: `go test ./...`
- Maintain or improve code coverage

### Project Structure

```
HomeGuard/
├── device/           # Device management module
├── listener/         # Listener implementations
├── wol/             # WOL core functionality
├── main.go          # Program entry point
└── ...
```

### Reporting Bugs

When reporting bugs, please include:

1. **Description**: Clear and concise description of the issue
2. **Steps to Reproduce**: Detailed reproduction steps
3. **Expected Behavior**: What you expected to happen
4. **Actual Behavior**: What actually happened
5. **Environment**:
   - Operating System
   - Go Version
   - HomeGuard Version
6. **Logs**: Relevant log output or error messages

### Feature Requests

When suggesting features, please include:

1. **Problem or Need**: Describe the problem you want to solve
2. **Proposed Solution**: Your idea
3. **Alternatives**: Other solutions you've considered
4. **Use Cases**: Specific application scenarios

### Code of Conduct

- Respect all contributors
- Be inclusive of different viewpoints and experiences
- Accept constructive criticism
- Focus on what's best for the community

---

<div align="center">
Thank you for contributing! / 感谢您的贡献！
</div>

