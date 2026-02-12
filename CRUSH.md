# Catwalk - AI提供商数据库

## 构建/测试命令

- `go run .` - 构建并运行主HTTP服务器，端口为:8080
- `go run ./cmd/{provider-name}` - 构建并运行CLI工具，用于更新`{provider-name}.json`文件
- `go test ./...` - 运行所有测试

## 代码风格指南

- 包注释：以"Package name provides/represents..."开头
- 导入顺序：标准库优先，然后是第三方库，最后是本地包
- 错误处理：使用`fmt.Errorf("message: %w", err)`包装错误
- 结构体标签：为可选字段使用带有omitempty的json标签
- 常量：将相关常量分组，并添加描述性注释
- 类型：为ID使用自定义类型（例如`InferenceProvider`, `Type`）
- 命名规则：未导出标识符使用驼峰式（camelCase），导出标识符使用帕斯卡式（PascalCase）
- 注释：使用`//nolint:directive`标记忽略特定的lint检查
- HTTP：始终设置超时，使用上下文，延迟关闭响应体
- JSON：使用`json.MarshalIndent`生成格式化输出，验证反序列化结果
- 文件权限：对敏感配置文件使用0o600权限

## 添加更多提供商命令

- 创建`./cmd/{provider-name}/main.go`文件
- 尝试使用提供商API获取可用模型列表。如果没有列出模型的端点，查找文档中的结构化文本格式。如果都不存在，则拒绝创建命令，并将其添加到`MANUAL_UPDATES.md`文件中。
- 将命令添加到`.github/workflows/update.yml`文件中

## 手动更新提供商

### Zai

对于`zai`，我们需要从`https://docs.z.ai/guides/overview/overview`获取模型列表和功能信息。

该页面不包含确切的`context_window`和`default_max_tokens`值。我们可以从`./internal/providers/configs/openrouter.json`文件中获取这些值。

