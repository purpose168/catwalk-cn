// Package main 提供了一个命令行工具，用于从 OpenRouter 获取模型
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

// Model 表示完整的模型配置。
type Model struct {
	ID              string       `json:"id"`
	CanonicalSlug   string       `json:"canonical_slug"`
	HuggingFaceID   string       `json:"hugging_face_id"`
	Name            string       `json:"name"`
	Created         int64        `json:"created"`
	Description     string       `json:"description"`
	ContextLength   int64        `json:"context_length"`
	Architecture    Architecture `json:"architecture"`
	Pricing         Pricing      `json:"pricing"`
	TopProvider     TopProvider  `json:"top_provider"`
	SupportedParams []string     `json:"supported_parameters"`
}

// Architecture 定义了模型的架构细节。
type Architecture struct {
	Modality         string   `json:"modality"`
	InputModalities  []string `json:"input_modalities"`
	OutputModalities []string `json:"output_modalities"`
	Tokenizer        string   `json:"tokenizer"`
	InstructType     *string  `json:"instruct_type"`
}

// Pricing 包含不同操作的定价信息。
type Pricing struct {
	Prompt            string `json:"prompt"`
	Completion        string `json:"completion"`
	Request           string `json:"request"`
	Image             string `json:"image"`
	WebSearch         string `json:"web_search"`
	InternalReasoning string `json:"internal_reasoning"`
	InputCacheRead    string `json:"input_cache_read"`
	InputCacheWrite   string `json:"input_cache_write"`
}

// TopProvider 描述了顶级提供程序的能力。
type TopProvider struct {
	ContextLength       int64  `json:"context_length"`
	MaxCompletionTokens *int64 `json:"max_completion_tokens"`
	IsModerated         bool   `json:"is_moderated"`
}

// Endpoint 表示模型的单个端点配置。
type Endpoint struct {
	Name                string   `json:"name"`
	ContextLength       int64    `json:"context_length"`
	Pricing             Pricing  `json:"pricing"`
	ProviderName        string   `json:"provider_name"`
	Tag                 string   `json:"tag"`
	Quantization        *string  `json:"quantization"`
	MaxCompletionTokens *int64   `json:"max_completion_tokens"`
	MaxPromptTokens     *int64   `json:"max_prompt_tokens"`
	SupportedParams     []string `json:"supported_parameters"`
	Status              int      `json:"status"`
	UptimeLast30m       float64  `json:"uptime_last_30m"`
}

// EndpointsResponse 是端点 API 的响应结构。
type EndpointsResponse struct {
	Data struct {
		ID          string     `json:"id"`
		Name        string     `json:"name"`
		Created     int64      `json:"created"`
		Description string     `json:"description"`
		Endpoints   []Endpoint `json:"endpoints"`
	} `json:"data"`
}

// ModelsResponse 是模型 API 的响应结构。
type ModelsResponse struct {
	Data []Model `json:"data"`
}

// ModelPricing 是模型的定价结构，详细说明了每百万令牌的输入和输出成本，包括缓存和非缓存。
type ModelPricing struct {
	CostPer1MIn        float64 `json:"cost_per_1m_in"`
	CostPer1MOut       float64 `json:"cost_per_1m_out"`
	CostPer1MInCached  float64 `json:"cost_per_1m_in_cached"`
	CostPer1MOutCached float64 `json:"cost_per_1m_out_cached"`
}

// getPricing 获取模型的定价信息
func getPricing(model Model) ModelPricing {
	pricing := ModelPricing{}
	costPrompt, err := strconv.ParseFloat(model.Pricing.Prompt, 64)
	if err != nil {
		costPrompt = 0.0
	}
	pricing.CostPer1MIn = costPrompt * 1_000_000
	costCompletion, err := strconv.ParseFloat(model.Pricing.Completion, 64)
	if err != nil {
		costCompletion = 0.0
	}
	pricing.CostPer1MOut = costCompletion * 1_000_000

	costPromptCached, err := strconv.ParseFloat(model.Pricing.InputCacheWrite, 64)
	if err != nil {
		costPromptCached = 0.0
	}
	pricing.CostPer1MInCached = costPromptCached * 1_000_000
	costCompletionCached, err := strconv.ParseFloat(model.Pricing.InputCacheRead, 64)
	if err != nil {
		costCompletionCached = 0.0
	}
	pricing.CostPer1MOutCached = costCompletionCached * 1_000_000
	return pricing
}

// fetchOpenRouterModels 获取 OpenRouter 模型列表
func fetchOpenRouterModels() (*ModelsResponse, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequestWithContext(
		context.Background(),
		"GET",
		"https://openrouter.ai/api/v1/models",
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

// fetchModelEndpoints 获取模型的端点列表
func fetchModelEndpoints(modelID string) (*EndpointsResponse, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	url := fmt.Sprintf("https://openrouter.ai/api/v1/models/%s/endpoints", modelID)
	req, _ := http.NewRequestWithContext(
		context.Background(),
		"GET",
		url,
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
	var er EndpointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		return nil, err //nolint:wrapcheck
	}
	return &er, nil
}

// selectBestEndpoint 选择最佳端点
func selectBestEndpoint(endpoints []Endpoint) *Endpoint {
	if len(endpoints) == 0 {
		return nil
	}

	var best *Endpoint
	for i := range endpoints {
		endpoint := &endpoints[i]
		// 跳过状态不佳或正常运行时间不足的端点
		if endpoint.Status < 0 || endpoint.UptimeLast30m < 90.0 {
			continue
		}

		if best == nil {
			best = endpoint
			continue
		}

		if isBetterEndpoint(endpoint, best) {
			best = endpoint
		}
	}

	// If no good endpoint found, return the first one as fallback
	if best == nil {
		best = &endpoints[0]
	}

	return best
}

// isBetterEndpoint 判断候选端点是否比当前端点更好
func isBetterEndpoint(candidate, current *Endpoint) bool {
	candidateHasTools := slices.Contains(candidate.SupportedParams, "tools")
	currentHasTools := slices.Contains(current.SupportedParams, "tools")

	// 优先选择支持工具的端点
	if candidateHasTools && !currentHasTools {
		return true
	}
	if !candidateHasTools && currentHasTools {
		return false
	}

	// 如果工具支持状态相同，则比较其他因素
	if candidate.ContextLength > current.ContextLength {
		return true
	}
	if candidate.ContextLength == current.ContextLength {
		return candidate.UptimeLast30m > current.UptimeLast30m
	}

	return false
}

// 用于生成 openrouter.json 配置文件。
func main() {
	modelsResp, err := fetchOpenRouterModels()
	if err != nil {
		log.Fatal("获取 OpenRouter 模型时出错:", err)
	}

	openRouterProvider := catwalk.Provider{
		Name:                "OpenRouter",
		ID:                  "openrouter",
		APIKey:              "$OPENROUTER_API_KEY",
		APIEndpoint:         "https://openrouter.ai/api/v1",
		Type:                catwalk.TypeOpenRouter,
		DefaultLargeModelID: "anthropic/claude-sonnet-4",
		DefaultSmallModelID: "anthropic/claude-3.5-haiku",
		Models:              []catwalk.Model{},
		DefaultHeaders: map[string]string{
			"HTTP-Referer": "https://charm.land",
			"X-Title":      "Crush",
		},
	}

	for _, model := range modelsResp.Data {
		if model.ContextLength < 20000 {
			continue
		}
		// 跳过非文本模型或不支持工具的模型
		if !slices.Contains(model.SupportedParams, "tools") ||
			!slices.Contains(model.Architecture.InputModalities, "text") ||
			!slices.Contains(model.Architecture.OutputModalities, "text") {
			continue
		}

		// 获取此模型的端点以获得最佳配置
		endpointsResp, err := fetchModelEndpoints(model.ID)
		if err != nil {
			fmt.Printf("警告：获取 %s 的端点失败: %v\n", model.ID, err)
			// 回退到使用原始模型数据
			pricing := getPricing(model)
			canReason := slices.Contains(model.SupportedParams, "reasoning")
			supportsImages := slices.Contains(model.Architecture.InputModalities, "image")

			var reasoningLevels []string
			var defaultReasoning string
			if canReason {
				reasoningLevels = []string{"low", "medium", "high"}
				defaultReasoning = "medium"
			}
			m := catwalk.Model{
				ID:                     model.ID,
				Name:                   model.Name,
				CostPer1MIn:            pricing.CostPer1MIn,
				CostPer1MOut:           pricing.CostPer1MOut,
				CostPer1MInCached:      pricing.CostPer1MInCached,
				CostPer1MOutCached:     pricing.CostPer1MOutCached,
				ContextWindow:          model.ContextLength,
				CanReason:              canReason,
				DefaultReasoningEffort: defaultReasoning,
				ReasoningLevels:        reasoningLevels,
				SupportsImages:         supportsImages,
			}
			if model.TopProvider.MaxCompletionTokens != nil {
				m.DefaultMaxTokens = *model.TopProvider.MaxCompletionTokens / 2
			} else {
				m.DefaultMaxTokens = model.ContextLength / 10
			}
			if model.TopProvider.ContextLength > 0 {
				m.ContextWindow = model.TopProvider.ContextLength
			}
			openRouterProvider.Models = append(openRouterProvider.Models, m)
			continue
		}

		// 选择最佳端点
		bestEndpoint := selectBestEndpoint(endpointsResp.Data.Endpoints)
		if bestEndpoint == nil {
			fmt.Printf("警告：未找到 %s 的合适端点\n", model.ID)
			continue
		}

		// 检查最佳端点是否支持工具
		if !slices.Contains(bestEndpoint.SupportedParams, "tools") {
			continue
		}

		// 使用最佳端点的配置
		pricing := ModelPricing{}
		costPrompt, err := strconv.ParseFloat(bestEndpoint.Pricing.Prompt, 64)
		if err != nil {
			costPrompt = 0.0
		}
		pricing.CostPer1MIn = costPrompt * 1_000_000
		costCompletion, err := strconv.ParseFloat(bestEndpoint.Pricing.Completion, 64)
		if err != nil {
			costCompletion = 0.0
		}
		pricing.CostPer1MOut = costCompletion * 1_000_000

		costPromptCached, err := strconv.ParseFloat(bestEndpoint.Pricing.InputCacheWrite, 64)
		if err != nil {
			costPromptCached = 0.0
		}
		pricing.CostPer1MInCached = costPromptCached * 1_000_000
		costCompletionCached, err := strconv.ParseFloat(bestEndpoint.Pricing.InputCacheRead, 64)
		if err != nil {
			costCompletionCached = 0.0
		}
		pricing.CostPer1MOutCached = costCompletionCached * 1_000_000

		canReason := slices.Contains(bestEndpoint.SupportedParams, "reasoning")
		supportsImages := slices.Contains(model.Architecture.InputModalities, "image")

		var reasoningLevels []string
		var defaultReasoning string
		if canReason {
			reasoningLevels = []string{"low", "medium", "high"}
			defaultReasoning = "medium"
		}
		m := catwalk.Model{
			ID:                     model.ID,
			Name:                   model.Name,
			CostPer1MIn:            pricing.CostPer1MIn,
			CostPer1MOut:           pricing.CostPer1MOut,
			CostPer1MInCached:      pricing.CostPer1MInCached,
			CostPer1MOutCached:     pricing.CostPer1MOutCached,
			ContextWindow:          bestEndpoint.ContextLength,
			CanReason:              canReason,
			DefaultReasoningEffort: defaultReasoning,
			ReasoningLevels:        reasoningLevels,
			SupportsImages:         supportsImages,
		}

		// 根据最佳端点设置最大令牌数
		if bestEndpoint.MaxCompletionTokens != nil {
			m.DefaultMaxTokens = *bestEndpoint.MaxCompletionTokens / 2
		} else {
			m.DefaultMaxTokens = bestEndpoint.ContextLength / 10
		}

		openRouterProvider.Models = append(openRouterProvider.Models, m)
		fmt.Printf("已添加模型 %s，上下文窗口为 %d，来自提供程序 %s\n",
			model.ID, bestEndpoint.ContextLength, bestEndpoint.ProviderName)
	}

	slices.SortFunc(openRouterProvider.Models, func(a catwalk.Model, b catwalk.Model) int {
		return strings.Compare(a.Name, b.Name)
	})

	// 将 JSON 保存到 internal/providers/config/openrouter.json
	data, err := json.MarshalIndent(openRouterProvider, "", "  ")
	if err != nil {
		log.Fatal("序列化 OpenRouter 提供程序时出错:", err)
	}
	// 写入文件
	if err := os.WriteFile("internal/providers/configs/openrouter.json", data, 0o600); err != nil {
		log.Fatal("写入 OpenRouter 提供程序配置时出错:", err)
	}
}
