// Package deprecated 用于提供旧版本的AI提供商配置
package deprecated

import "github.com/purpose168/catwalk-cn/pkg/catwalk"

// Provider 表示AI提供商的配置信息
type Provider struct {
	Name                string                    `json:"name"`                             // 提供商名称
	ID                  catwalk.InferenceProvider `json:"id"`                               // 推理提供商ID
	APIKey              string                    `json:"api_key,omitempty"`                // API密钥
	APIEndpoint         string                    `json:"api_endpoint,omitempty"`           // API端点地址
	Type                catwalk.Type              `json:"type,omitempty"`                   // 类型
	DefaultLargeModelID string                    `json:"default_large_model_id,omitempty"` // 默认大模型ID
	DefaultSmallModelID string                    `json:"default_small_model_id,omitempty"` // 默认小模型ID
	Models              []Model                   `json:"models,omitempty"`                 // 模型列表
	DefaultHeaders      map[string]string         `json:"default_headers,omitempty"`        // 默认请求头
}

// Model 表示AI模型的配置信息
type Model struct {
	ID                     string  `json:"id"`                                 // 模型ID
	Name                   string  `json:"name"`                               // 模型名称
	CostPer1MIn            float64 `json:"cost_per_1m_in"`                     // 每百万输入token的成本
	CostPer1MOut           float64 `json:"cost_per_1m_out"`                    // 每百万输出token的成本
	CostPer1MInCached      float64 `json:"cost_per_1m_in_cached"`              // 每百万缓存输入token的成本
	CostPer1MOutCached     float64 `json:"cost_per_1m_out_cached"`             // 每百万缓存输出token的成本
	ContextWindow          int64   `json:"context_window"`                     // 上下文窗口大小
	DefaultMaxTokens       int64   `json:"default_max_tokens"`                 // 默认最大token数
	CanReason              bool    `json:"can_reason"`                         // 是否支持推理
	HasReasoningEffort     bool    `json:"has_reasoning_efforts"`              // 是否有推理努力参数
	DefaultReasoningEffort string  `json:"default_reasoning_effort,omitempty"` // 默认推理努力值
	SupportsImages         bool    `json:"supports_attachments"`               // 是否支持图片附件
}
