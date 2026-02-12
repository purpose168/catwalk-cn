// Package main 提供了一个命令行工具，用于从 Synthetic 获取模型
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

// Model 表示来自 Synthetic API 的模型。
type Model struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	InputModalities   []string `json:"input_modalities"`
	OutputModalities  []string `json:"output_modalities"`
	ContextLength     int64    `json:"context_length"`
	MaxOutputLength   int64    `json:"max_output_length,omitempty"`
	Pricing           Pricing  `json:"pricing"`
	SupportedFeatures []string `json:"supported_features,omitempty"`
}

// Pricing 包含不同操作的定价信息。
type Pricing struct {
	Prompt           string `json:"prompt"`
	Completion       string `json:"completion"`
	Image            string `json:"image"`
	Request          string `json:"request"`
	InputCacheReads  string `json:"input_cache_reads"`
	InputCacheWrites string `json:"input_cache_writes"`
}

// ModelsResponse 是 Synthetic 模型 API 的响应结构。
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

// parsePrice 从 Synthetic 的价格格式（例如 "$0.00000055"）中提取浮点数。
func parsePrice(s string) float64 {
	s = strings.TrimPrefix(s, "$")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return v
}

// getPricing 获取模型的定价信息
func getPricing(model Model) ModelPricing {
	return ModelPricing{
		CostPer1MIn:        parsePrice(model.Pricing.Prompt) * 1_000_000,
		CostPer1MOut:       parsePrice(model.Pricing.Completion) * 1_000_000,
		CostPer1MInCached:  parsePrice(model.Pricing.InputCacheReads) * 1_000_000,
		CostPer1MOutCached: parsePrice(model.Pricing.InputCacheReads) * 1_000_000,
	}
}

// applyModelOverrides 为 Synthetic 省略元数据的模型设置 supported_features。
// TODO: 当他们添加缺失的元数据时，删除此函数。
func applyModelOverrides(model *Model) {
	switch {
	// 所有 llama 模型都支持工具，但目前都不支持推理
	case strings.HasPrefix(model.ID, "hf:meta-llama/Llama-"):
		model.SupportedFeatures = []string{"tools"}

	case strings.HasPrefix(model.ID, "hf:deepseek-ai/DeepSeek-R1"):
		model.SupportedFeatures = []string{"tools", "reasoning"}

	case strings.HasPrefix(model.ID, "hf:deepseek-ai/DeepSeek-V3.1"):
		model.SupportedFeatures = []string{"tools", "reasoning"}

	case strings.HasPrefix(model.ID, "hf:deepseek-ai/DeepSeek-V3.2"):
		model.SupportedFeatures = []string{"tools", "reasoning"}

	case strings.HasPrefix(model.ID, "hf:deepseek-ai/DeepSeek-V3"):
		model.SupportedFeatures = []string{"tools"}

	case strings.HasPrefix(model.ID, "hf:Qwen/Qwen3-235B-A22B-Thinking"):
		model.SupportedFeatures = []string{"tools", "reasoning"}

	case strings.HasPrefix(model.ID, "hf:Qwen/Qwen3-235B-A22B-Instruct"):
		model.SupportedFeatures = []string{"tools", "reasoning"}

	// 其余的 Qwen3 模型不支持推理，但支持工具
	case strings.HasPrefix(model.ID, "hf:Qwen/Qwen3"):
		model.SupportedFeatures = []string{"tools"}

	// 已经有正确的元数据，但下面的 k2 匹配器会覆盖它以省略推理
	case strings.HasPrefix(model.ID, "hf:moonshotai/Kimi-K2-Thinking"):
		model.SupportedFeatures = []string{"tools", "reasoning"}

	case strings.HasPrefix(model.ID, "hf:moonshotai/Kimi-K2.5"):
		model.SupportedFeatures = []string{"tools", "reasoning"}

	case strings.HasPrefix(model.ID, "hf:moonshotai/Kimi-K2"):
		model.SupportedFeatures = []string{"tools"}

	case strings.HasPrefix(model.ID, "hf:zai-org/GLM-4.5"):
		model.SupportedFeatures = []string{"tools"}

	case strings.HasPrefix(model.ID, "hf:openai/gpt-oss"):
		model.SupportedFeatures = []string{"tools", "reasoning"}

	case strings.HasPrefix(model.ID, "hf:MiniMaxAI/MiniMax-M2.1"):
		model.SupportedFeatures = []string{"tools", "reasoning"}
	}
}

// fetchSyntheticModels 获取 Synthetic 模型列表
func fetchSyntheticModels(apiEndpoint string) (*ModelsResponse, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequestWithContext(context.Background(), "GET", apiEndpoint+"/models", nil)
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

// 用于生成 synthetic.json 配置文件。
func main() {
	syntheticProvider := catwalk.Provider{
		Name:                "Synthetic",
		ID:                  "synthetic",
		APIKey:              "$SYNTHETIC_API_KEY",
		APIEndpoint:         "https://api.synthetic.new/openai/v1",
		Type:                catwalk.TypeOpenAICompat,
		DefaultLargeModelID: "hf:zai-org/GLM-4.7",
		DefaultSmallModelID: "hf:deepseek-ai/DeepSeek-V3.1-Terminus",
		Models:              []catwalk.Model{},
	}

	modelsResp, err := fetchSyntheticModels(syntheticProvider.APIEndpoint)
	if err != nil {
		log.Fatal("获取 Synthetic 模型时出错:", err)
	}

	// 为缺少 supported_features 元数据的模型应用覆盖
	for i := range modelsResp.Data {
		applyModelOverrides(&modelsResp.Data[i])
	}

	for _, model := range modelsResp.Data {
		// 跳过上下文窗口较小的模型
		if model.ContextLength < 20000 {
			continue
		}

		// 跳过非文本模型
		if !slices.Contains(model.InputModalities, "text") ||
			!slices.Contains(model.OutputModalities, "text") {
			continue
		}

		// 确保它们支持工具
		supportsTools := slices.Contains(model.SupportedFeatures, "tools")
		if !supportsTools {
			continue
		}

		pricing := getPricing(model)
		supportsImages := slices.Contains(model.InputModalities, "image")

		// Check if model supports reasoning
		canReason := slices.Contains(model.SupportedFeatures, "reasoning")
		var reasoningLevels []string
		var defaultReasoning string
		if canReason {
			reasoningLevels = []string{"low", "medium", "high"}
			defaultReasoning = "medium"
		}

		// 去除第一个 / 之前的所有内容，以获得更简洁的名称
		modelName := model.Name
		if idx := strings.Index(model.Name, "/"); idx != -1 {
			modelName = model.Name[idx+1:]
		}
		// 将连字符替换为空格
		modelName = strings.ReplaceAll(modelName, "-", " ")

		m := catwalk.Model{
			ID:                     model.ID,
			Name:                   modelName,
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

		// 如果 max_output_length 可用，则基于它设置最大令牌数，但上限为上下文长度的 15%
		maxFromOutput := model.MaxOutputLength / 2
		maxAt15Pct := (model.ContextLength * 15) / 100
		if model.MaxOutputLength > 0 && maxFromOutput <= maxAt15Pct {
			m.DefaultMaxTokens = maxFromOutput
		} else {
			m.DefaultMaxTokens = model.ContextLength / 10
		}

		syntheticProvider.Models = append(syntheticProvider.Models, m)
		fmt.Printf("已添加模型 %s，上下文窗口为 %d\n",
			model.ID, model.ContextLength)
	}

	slices.SortFunc(syntheticProvider.Models, func(a catwalk.Model, b catwalk.Model) int {
		return strings.Compare(a.Name, b.Name)
	})

	// 将 JSON 保存到 internal/providers/configs/synthetic.json
	data, err := json.MarshalIndent(syntheticProvider, "", "  ")
	if err != nil {
		log.Fatal("序列化 Synthetic 提供程序时出错:", err)
	}

	if err := os.WriteFile("internal/providers/configs/synthetic.json", data, 0o600); err != nil {
		log.Fatal("写入 Synthetic 提供程序配置时出错:", err)
	}

	fmt.Printf("已生成 synthetic.json，包含 %d 个模型\n", len(syntheticProvider.Models))
}
