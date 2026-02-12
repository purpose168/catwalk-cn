// Package main 提供了一个命令行工具，用于从 Vercel AI Gateway 获取模型
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
	"strconv"
	"strings"
	"time"

	"github.com/purpose168/catwalk-cn/pkg/catwalk"
)

// Model 表示来自 Vercel API 的模型。
type Model struct {
	ID            string   `json:"id"`
	Object        string   `json:"object"`
	Created       int64    `json:"created"`
	OwnedBy       string   `json:"owned_by"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	ContextWindow int64    `json:"context_window"`
	MaxTokens     int64    `json:"max_tokens"`
	Type          string   `json:"type"`
	Tags          []string `json:"tags"`
	Pricing       Pricing  `json:"pricing"`
}

// Pricing 包含模型的定价信息。
type Pricing struct {
	Input           string `json:"input,omitempty"`
	Output          string `json:"output,omitempty"`
	InputCacheRead  string `json:"input_cache_read,omitempty"`
	InputCacheWrite string `json:"input_cache_write,omitempty"`
	WebSearch       string `json:"web_search,omitempty"`
	Image           string `json:"image,omitempty"`
}

// ModelsResponse 是 Vercel 模型 API 的响应结构。
type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// fetchVercelModels 获取 Vercel 模型列表
func fetchVercelModels() (*ModelsResponse, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequestWithContext(
		context.Background(),
		"GET",
		"https://ai-gateway.vercel.sh/v1/models",
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

func main() {
	modelsResp, err := fetchVercelModels()
	if err != nil {
		log.Fatal("获取 Vercel 模型时出错:", err)
	}

	vercelProvider := catwalk.Provider{
		Name:                "Vercel",
		ID:                  catwalk.InferenceProviderVercel,
		APIKey:              "$VERCEL_API_KEY",
		APIEndpoint:         "https://ai-gateway.vercel.sh/v1",
		Type:                catwalk.TypeVercel,
		DefaultLargeModelID: "anthropic/claude-sonnet-4",
		DefaultSmallModelID: "anthropic/claude-haiku-4.5",
		Models:              []catwalk.Model{},
		DefaultHeaders: map[string]string{
			"HTTP-Referer": "https://charm.land",
			"X-Title":      "Crush",
		},
	}

	for _, model := range modelsResp.Data {
		// 只包含语言模型，跳过嵌入和图像模型
		if model.Type != "language" {
			continue
		}

		// 跳过不支持工具的模型
		if !slices.Contains(model.Tags, "tool-use") {
			continue
		}

		// 解析定价
		costPer1MIn := 0.0
		costPer1MOut := 0.0
		costPer1MInCached := 0.0
		costPer1MOutCached := 0.0

		if model.Pricing.Input != "" {
			costPrompt, err := strconv.ParseFloat(model.Pricing.Input, 64)
			if err == nil {
				costPer1MIn = costPrompt * 1_000_000
			}
		}

		if model.Pricing.Output != "" {
			costCompletion, err := strconv.ParseFloat(model.Pricing.Output, 64)
			if err == nil {
				costPer1MOut = costCompletion * 1_000_000
			}
		}

		if model.Pricing.InputCacheRead != "" {
			costCached, err := strconv.ParseFloat(model.Pricing.InputCacheRead, 64)
			if err == nil {
				costPer1MInCached = costCached * 1_000_000
			}
		}

		if model.Pricing.InputCacheWrite != "" {
			costCacheWrite, err := strconv.ParseFloat(model.Pricing.InputCacheWrite, 64)
			if err == nil {
				costPer1MOutCached = costCacheWrite * 1_000_000
			}
		}

		// 检查模型是否支持推理
		canReason := slices.Contains(model.Tags, "reasoning")

		var reasoningLevels []string
		var defaultReasoning string
		if canReason {
			// 大多数提供商支持的基础推理级别
			reasoningLevels = []string{"low", "medium", "high"}
			// Anthropic 模型支持扩展的 Vercel 推理级别
			if strings.HasPrefix(model.ID, "anthropic/") {
				reasoningLevels = []string{"none", "minimal", "low", "medium", "high", "xhigh"}
			}
			defaultReasoning = "medium"
		}

		// 检查模型是否支持图像
		supportsImages := slices.Contains(model.Tags, "vision")

		// 计算默认最大令牌数
		defaultMaxTokens := model.MaxTokens
		if defaultMaxTokens == 0 {
			defaultMaxTokens = model.ContextWindow / 10
		}
		if defaultMaxTokens > 8000 {
			defaultMaxTokens = 8000
		}

		m := catwalk.Model{
			ID:                     model.ID,
			Name:                   model.Name,
			CostPer1MIn:            costPer1MIn,
			CostPer1MOut:           costPer1MOut,
			CostPer1MInCached:      costPer1MInCached,
			CostPer1MOutCached:     costPer1MOutCached,
			ContextWindow:          model.ContextWindow,
			DefaultMaxTokens:       defaultMaxTokens,
			CanReason:              canReason,
			ReasoningLevels:        reasoningLevels,
			DefaultReasoningEffort: defaultReasoning,
			SupportsImages:         supportsImages,
		}

		vercelProvider.Models = append(vercelProvider.Models, m)
		fmt.Printf("已添加模型 %s，上下文窗口为 %d\n", model.ID, model.ContextWindow)
	}

	slices.SortFunc(vercelProvider.Models, func(a catwalk.Model, b catwalk.Model) int {
		return strings.Compare(a.Name, b.Name)
	})

	// 将 JSON 保存到 internal/providers/configs/vercel.json
	data, err := json.MarshalIndent(vercelProvider, "", "  ")
	if err != nil {
		log.Fatal("序列化 Vercel 提供程序时出错:", err)
	}

	if err := os.WriteFile("internal/providers/configs/vercel.json", data, 0o600); err != nil {
		log.Fatal("写入 Vercel 提供程序配置时出错:", err)
	}

	fmt.Printf("已生成 vercel.json，包含 %d 个模型\n", len(vercelProvider.Models))
}
