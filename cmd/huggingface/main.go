// Package main 提供了一个命令行工具，用于从 Hugging Face Router 获取模型
// 并为提供程序生成配置文件。
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/purpose168/catwalk-cn/pkg/catwalk"
)

// SupportedProviders 定义了我们想要支持的提供程序。
// 可以通过添加或删除此切片中的提供程序来控制包含哪些提供程序。
var SupportedProviders = []string{
	// "together", // 存在多个问题
	"fireworks-ai",
	//"nebius",
	// "novita", // 使用报告不正确
	"groq",
	"cerebras",
	// "hyperbolic",
	// "nscale",
	// "sambanova",
	// "cohere",
	"hf-inference",
}

// Model 表示来自 Hugging Face Router API 的模型。
type Model struct {
	ID        string     `json:"id"`
	Object    string     `json:"object"`
	Created   int64      `json:"created"`
	OwnedBy   string     `json:"owned_by"`
	Providers []Provider `json:"providers"`
}

// Provider 表示模型的提供程序配置。
type Provider struct {
	Provider                 string   `json:"provider"`
	Status                   string   `json:"status"`
	ContextLength            int64    `json:"context_length,omitempty"`
	Pricing                  *Pricing `json:"pricing,omitempty"`
	SupportsTools            bool     `json:"supports_tools"`
	SupportsStructuredOutput bool     `json:"supports_structured_output"`
}

// Pricing 包含提供程序的定价信息。
type Pricing struct {
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
}

// ModelsResponse 是 Hugging Face Router 模型 API 的响应结构。
type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// fetchHuggingFaceModels 获取 Hugging Face 模型列表
func fetchHuggingFaceModels() (*ModelsResponse, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequestWithContext(
		context.Background(),
		"GET",
		"https://router.huggingface.co/v1/models",
		nil,
	)
	req.Header.Set("User-Agent", "Crush-Client/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("状态码 %d: %s", resp.StatusCode, body)
	}
	var mr ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&mr); err != nil {
		return nil, err //nolint:wrapcheck
	}
	return &mr, nil
}

// findContextWindow 查找给定模型的任何提供程序的上下文窗口大小。
func findContextWindow(model Model) int64 {
	for _, provider := range model.Providers {
		if provider.ContextLength > 0 {
			return provider.ContextLength
		}
	}
	return 0
}

// 警告：请勿使用
// 目前我们只使用模型的子集列表。
func main() {
	modelsResp, err := fetchHuggingFaceModels()
	if err != nil {
		log.Fatal("获取 Hugging Face 模型时出错:", err)
	}

	hfProvider := catwalk.Provider{
		Name:                "Hugging Face",
		ID:                  catwalk.InferenceProviderHuggingFace,
		APIKey:              "$HF_TOKEN",
		APIEndpoint:         "https://router.huggingface.co/v1",
		Type:                catwalk.TypeOpenAICompat,
		DefaultLargeModelID: "moonshotai/Kimi-K2-Instruct-0905:groq",
		DefaultSmallModelID: "openai/gpt-oss-20b",
		Models:              []catwalk.Model{},
		DefaultHeaders: map[string]string{
			"HTTP-Referer": "https://charm.land",
			"X-Title":      "Crush",
		},
	}

	for _, model := range modelsResp.Data {
		// 查找此模型的任何提供程序的上下文窗口
		fallbackContextLength := findContextWindow(model)
		if fallbackContextLength == 0 {
			fmt.Printf("跳过模型 %s - 未在任何提供程序中找到上下文窗口\n", model.ID)
			continue
		}

		for _, provider := range model.Providers {
			// 跳过不支持的提供程序
			if !slices.Contains(SupportedProviders, provider.Provider) {
				continue
			}

			// 跳过不支持工具的提供程序
			if !provider.SupportsTools {
				continue
			}

			// 跳过非活跃的提供程序
			if provider.Status != "live" {
				continue
			}

			// 创建带有提供程序特定 ID 和名称的模型
			modelID := fmt.Sprintf("%s:%s", model.ID, provider.Provider)
			modelName := fmt.Sprintf("%s (%s)", model.ID, provider.Provider)

			// 使用提供程序的上下文长度，如果不可用则使用回退值
			contextLength := provider.ContextLength
			if contextLength == 0 {
				contextLength = fallbackContextLength
			}

			// 计算定价（从每令牌转换为每百万令牌）
			var costPer1MIn, costPer1MOut float64
			if provider.Pricing != nil {
				costPer1MIn = provider.Pricing.Input
				costPer1MOut = provider.Pricing.Output
			}

			// 设置默认最大令牌数（保守估计）
			defaultMaxTokens := min(contextLength/4, 8192)

			m := catwalk.Model{
				ID:                 modelID,
				Name:               modelName,
				CostPer1MIn:        costPer1MIn,
				CostPer1MOut:       costPer1MOut,
				CostPer1MInCached:  0, // 未由 HF Router 提供
				CostPer1MOutCached: 0, // 未由 HF Router 提供
				ContextWindow:      contextLength,
				DefaultMaxTokens:   defaultMaxTokens,
				CanReason:          false, // 未由 HF Router 提供
				SupportsImages:     false, // 未由 HF Router 提供
			}

			hfProvider.Models = append(hfProvider.Models, m)
			fmt.Printf("已添加模型 %s，上下文窗口为 %d，来自提供程序 %s\n",
				modelID, contextLength, provider.Provider)
		}
	}

	slices.SortFunc(hfProvider.Models, func(a catwalk.Model, b catwalk.Model) int {
		return strings.Compare(a.Name, b.Name)
	})

	// 将 JSON 保存到 internal/providers/configs/huggingface.json
	data, err := json.MarshalIndent(hfProvider, "", "  ")
	if err != nil {
		log.Fatal("序列化 Hugging Face 提供程序时出错:", err)
	}

	if err := os.WriteFile("internal/providers/configs/huggingface.json", data, 0o600); err != nil {
		log.Fatal("写入 Hugging Face 提供程序配置时出错:", err)
	}

	fmt.Printf("已生成 huggingface.json，包含 %d 个模型\n", len(hfProvider.Models))
}
