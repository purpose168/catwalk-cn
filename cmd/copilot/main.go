// Package main 实现了一个工具，用于获取 GitHub Copilot 模型并生成 Catwalk 提供程序配置。
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/purpose168/catwalk-cn/pkg/catwalk"
)

// Response 表示 API 返回的响应结构
type Response struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// APITokenResponse 表示 API 令牌响应结构
type APITokenResponse struct {
	Token     string                    `json:"token"`
	ExpiresAt int64                     `json:"expires_at"`
	Endpoints APITokenResponseEndpoints `json:"endpoints"`
}

// APITokenResponseEndpoints 表示 API 令牌响应中的端点信息
type APITokenResponseEndpoints struct {
	API string `json:"api"`
}

// APIToken 表示 API 令牌结构
type APIToken struct {
	APIKey      string
	ExpiresAt   time.Time
	APIEndpoint string
}

// Model 表示模型信息结构
type Model struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	Version            string     `json:"version"`
	Vendor             string     `json:"vendor"`
	Preview            bool       `json:"preview"`
	ModelPickerEnabled bool       `json:"model_picker_enabled"`
	Capabilities       Capability `json:"capabilities"`
	Policy             *Policy    `json:"policy,omitempty"`
}

// Capability 表示模型能力结构
type Capability struct {
	Family    string   `json:"family"`
	Type      string   `json:"type"`
	Tokenizer string   `json:"tokenizer"`
	Limits    Limits   `json:"limits"`
	Supports  Supports `json:"supports"`
}

// Limits 表示模型限制结构
type Limits struct {
	MaxContextWindowTokens int `json:"max_context_window_tokens,omitempty"`
	MaxOutputTokens        int `json:"max_output_tokens,omitempty"`
	MaxPromptTokens        int `json:"max_prompt_tokens,omitempty"`
}

// Supports 表示模型支持的功能结构
type Supports struct {
	ToolCalls         bool `json:"tool_calls,omitempty"`
	ParallelToolCalls bool `json:"parallel_tool_calls,omitempty"`
	Vision            bool `json:"vision,omitempty"`
}

// Policy 表示模型策略结构
type Policy struct {
	State string `json:"state"`
	Terms string `json:"terms"`
}

var versionedModelRegexp = regexp.MustCompile(`-\d{4}-\d{2}-\d{2}$`)

// main 是程序的入口函数
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run 是程序的主要执行函数
func run() error {
	copilotModels, err := fetchCopilotModels()
	if err != nil {
		return err
	}

	// NOTE(@andreynering): 排除版本化模型，只保留每个模型的主版本。
	copilotModels = slices.DeleteFunc(copilotModels, func(m Model) bool {
		return m.ID != m.Version || versionedModelRegexp.MatchString(m.ID) || strings.Contains(m.ID, "embedding")
	})

	catwalkModels := modelsToCatwalk(copilotModels)
	slices.SortStableFunc(catwalkModels, func(a, b catwalk.Model) int {
		return strings.Compare(a.ID, b.ID)
	})

	provider := catwalk.Provider{
		ID:                  catwalk.InferenceProviderCopilot,
		Name:                "GitHub Copilot",
		Models:              catwalkModels,
		APIEndpoint:         "https://api.githubcopilot.com",
		Type:                catwalk.TypeOpenAICompat,
		DefaultLargeModelID: "claude-opus-4.6",
		DefaultSmallModelID: "claude-haiku-4.5",
	}
	data, err := json.MarshalIndent(provider, "", "  ")
	if err != nil {
		return fmt.Errorf("无法序列化 JSON: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile("internal/providers/configs/copilot.json", data, 0o600); err != nil {
		return fmt.Errorf("无法写入 copilot.json: %w", err)
	}
	return nil
}

// fetchCopilotModels 用于获取 Copilot 模型列表
func fetchCopilotModels() ([]Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	oauthToken := copilotToken()
	if oauthToken == "" {
		return nil, fmt.Errorf("没有可用的 OAuth 令牌")
	}

	// 步骤1：从令牌端点获取 API 令牌
	tokenURL := "https://api.github.com/copilot_internal/v2/token" //nolint:gosec
	tokenReq, err := http.NewRequestWithContext(ctx, "GET", tokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("无法创建令牌请求: %w", err)
	}
	tokenReq.Header.Set("Accept", "application/json")
	tokenReq.Header.Set("Authorization", fmt.Sprintf("token %s", oauthToken))

	// 使用已批准的集成 ID 绕过客户端检查
	tokenReq.Header.Set("Copilot-Integration-Id", "vscode-chat")
	tokenReq.Header.Set("User-Agent", "GitHubCopilotChat/0.1")

	client := &http.Client{}
	tokenResp, err := client.Do(tokenReq)
	if err != nil {
		return nil, fmt.Errorf("无法发送令牌请求: %w", err)
	}
	defer tokenResp.Body.Close() //nolint:errcheck

	tokenBody, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return nil, fmt.Errorf("无法读取令牌响应体: %w", err)
	}

	if tokenResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("令牌端点返回意外状态码: %d", tokenResp.StatusCode)
	}

	var tokenData APITokenResponse
	if err := json.Unmarshal(tokenBody, &tokenData); err != nil {
		return nil, fmt.Errorf("无法解析令牌响应: %w", err)
	}

	// 转换为 APIToken 类型
	expiresAt := time.Unix(tokenData.ExpiresAt, 0)
	apiToken := APIToken{
		APIKey:      tokenData.Token,
		ExpiresAt:   expiresAt,
		APIEndpoint: tokenData.Endpoints.API,
	}

	// 步骤2：使用令牌中的动态端点获取模型列表
	modelsURL := apiToken.APIEndpoint + "/models"
	modelsReq, err := http.NewRequestWithContext(ctx, "GET", modelsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("无法创建模型请求: %w", err)
	}
	modelsReq.Header.Set("Accept", "application/json")
	modelsReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken.APIKey))
	modelsReq.Header.Set("Copilot-Integration-Id", "vscode-chat")
	modelsReq.Header.Set("User-Agent", "GitHubCopilotChat/0.1")

	modelsResp, err := client.Do(modelsReq)
	if err != nil {
		return nil, fmt.Errorf("无法发送模型请求: %w", err)
	}
	defer modelsResp.Body.Close() //nolint:errcheck

	modelsBody, err := io.ReadAll(modelsResp.Body)
	if err != nil {
		return nil, fmt.Errorf("无法读取模型响应体: %w", err)
	}

	if modelsResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("模型端点返回意外状态码: %d", modelsResp.StatusCode)
	}

	// 用于调试
	_ = os.MkdirAll("tmp", 0o700)
	_ = os.WriteFile("tmp/copilot-response.json", modelsBody, 0o600)

	var data Response
	if err := json.Unmarshal(modelsBody, &data); err != nil {
		return nil, fmt.Errorf("无法解析 JSON: %w", err)
	}
	return data.Data, nil
}

// modelsToCatwalk 将 Copilot 模型转换为 Catwalk 模型
func modelsToCatwalk(m []Model) []catwalk.Model {
	models := make([]catwalk.Model, 0, len(m))
	for _, model := range m {
		models = append(models, modelToCatwalk(model))
	}
	return models
}

// modelToCatwalk 将单个 Copilot 模型转换为 Catwalk 模型
func modelToCatwalk(m Model) catwalk.Model {
	return catwalk.Model{
		ID:               m.ID,
		Name:             m.Name,
		DefaultMaxTokens: int64(m.Capabilities.Limits.MaxOutputTokens),
		ContextWindow:    int64(m.Capabilities.Limits.MaxContextWindowTokens),
		SupportsImages:   m.Capabilities.Supports.Vision,
	}
}

// copilotToken 获取 Copilot 令牌
func copilotToken() string {
	if token := os.Getenv("COPILOT_TOKEN"); token != "" {
		return token
	}
	return tokenFromDisk()
}

// tokenFromDisk 从磁盘读取令牌
func tokenFromDisk() string {
	data, err := os.ReadFile(tokenFilePath())
	if err != nil {
		return ""
	}
	var content map[string]struct {
		User        string `json:"user"`
		OAuthToken  string `json:"oauth_token"`
		GitHubAppID string `json:"githubAppId"`
	}
	if err := json.Unmarshal(data, &content); err != nil {
		return ""
	}
	if app, ok := content["github.com:Iv1.b507a08c87ecfe98"]; ok {
		return app.OAuthToken
	}
	return ""
}

// tokenFilePath 获取令牌文件路径
func tokenFilePath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "github-copilot/apps.json")
	default:
		return filepath.Join(os.Getenv("HOME"), ".config/github-copilot/apps.json")
	}
}
