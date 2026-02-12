package catwalk

// Type 表示AI提供商的类型
type Type string

// 所有支持的AI提供商类型
const (
	TypeOpenAI       Type = "openai"        // OpenAI类型
	TypeOpenAICompat Type = "openai-compat" // OpenAI兼容类型
	TypeOpenRouter   Type = "openrouter"    // OpenRouter类型
	TypeVercel       Type = "vercel"        // Vercel类型
	TypeAnthropic    Type = "anthropic"     // Anthropic类型
	TypeGoogle       Type = "google"        // Google类型
	TypeAzure        Type = "azure"         // Azure类型
	TypeBedrock      Type = "bedrock"       // Bedrock类型
	TypeVertexAI     Type = "google-vertex" // Google Vertex AI类型
)

// InferenceProvider 表示推理提供商的标识符
type InferenceProvider string

// 系统支持的所有推理提供商
const (
	InferenceProviderOpenAI      InferenceProvider = "openai"      // OpenAI
	InferenceProviderAnthropic   InferenceProvider = "anthropic"   // Anthropic
	InferenceProviderSynthetic   InferenceProvider = "synthetic"   // Synthetic
	InferenceProviderGemini      InferenceProvider = "gemini"      // Gemini
	InferenceProviderAzure       InferenceProvider = "azure"       // Azure
	InferenceProviderBedrock     InferenceProvider = "bedrock"     // Bedrock
	InferenceProviderVertexAI    InferenceProvider = "vertexai"    // VertexAI
	InferenceProviderXAI         InferenceProvider = "xai"         // xAI
	InferenceProviderZAI         InferenceProvider = "zai"         // zAI
	InferenceProviderGROQ        InferenceProvider = "groq"        // GROQ
	InferenceProviderOpenRouter  InferenceProvider = "openrouter"  // OpenRouter
	InferenceProviderCerebras    InferenceProvider = "cerebras"    // Cerebras
	InferenceProviderVenice      InferenceProvider = "venice"      // Venice
	InferenceProviderChutes      InferenceProvider = "chutes"      // Chutes
	InferenceProviderHuggingFace InferenceProvider = "huggingface" // HuggingFace
	InferenceAIHubMix            InferenceProvider = "aihubmix"    // AIHubMix
	InferenceKimiCoding          InferenceProvider = "kimi-coding" // Kimi Coding
	InferenceProviderCopilot     InferenceProvider = "copilot"     // Copilot
	InferenceProviderVercel      InferenceProvider = "vercel"      // Vercel
	InferenceProviderMiniMax     InferenceProvider = "minimax"     // MiniMax
)

// Provider 表示AI提供商的配置信息
type Provider struct {
	Name                string            `json:"name"`                             // 提供商名称
	ID                  InferenceProvider `json:"id"`                               // 推理提供商ID
	APIKey              string            `json:"api_key,omitempty"`                // API密钥
	APIEndpoint         string            `json:"api_endpoint,omitempty"`           // API端点地址
	Type                Type              `json:"type,omitempty"`                   // 类型
	DefaultLargeModelID string            `json:"default_large_model_id,omitempty"` // 默认大模型ID
	DefaultSmallModelID string            `json:"default_small_model_id,omitempty"` // 默认小模型ID
	Models              []Model           `json:"models,omitempty"`                 // 模型列表
	DefaultHeaders      map[string]string `json:"default_headers,omitempty"`        // 默认请求头
}

// ModelOptions 存储模型的额外配置选项
type ModelOptions struct {
	Temperature      *float64       `json:"temperature,omitempty"`       // 温度参数
	TopP             *float64       `json:"top_p,omitempty"`             // TopP参数
	TopK             *int64         `json:"top_k,omitempty"`             // TopK参数
	FrequencyPenalty *float64       `json:"frequency_penalty,omitempty"` // 频率惩罚
	PresencePenalty  *float64       `json:"presence_penalty,omitempty"`  // 存在惩罚
	ProviderOptions  map[string]any `json:"provider_options,omitempty"`  // 提供商特定选项
}

// Model 表示AI模型的配置信息
type Model struct {
	ID                     string       `json:"id"`                                 // 模型ID
	Name                   string       `json:"name"`                               // 模型名称
	CostPer1MIn            float64      `json:"cost_per_1m_in"`                     // 每百万输入token的成本
	CostPer1MOut           float64      `json:"cost_per_1m_out"`                    // 每百万输出token的成本
	CostPer1MInCached      float64      `json:"cost_per_1m_in_cached"`              // 每百万缓存输入token的成本
	CostPer1MOutCached     float64      `json:"cost_per_1m_out_cached"`             // 每百万缓存输出token的成本
	ContextWindow          int64        `json:"context_window"`                     // 上下文窗口大小
	DefaultMaxTokens       int64        `json:"default_max_tokens"`                 // 默认最大token数
	CanReason              bool         `json:"can_reason"`                         // 是否支持推理
	ReasoningLevels        []string     `json:"reasoning_levels,omitempty"`         // 推理级别列表
	DefaultReasoningEffort string       `json:"default_reasoning_effort,omitempty"` // 默认推理努力值
	SupportsImages         bool         `json:"supports_attachments"`               // 是否支持图片附件
	Options                ModelOptions `json:"options"`                            // 模型配置选项
}

// KnownProviders 返回所有已知的推理提供商
func KnownProviders() []InferenceProvider {
	return []InferenceProvider{
		InferenceProviderOpenAI,
		InferenceProviderSynthetic,
		InferenceProviderAnthropic,
		InferenceProviderGemini,
		InferenceProviderAzure,
		InferenceProviderBedrock,
		InferenceProviderVertexAI,
		InferenceProviderXAI,
		InferenceProviderZAI,
		InferenceProviderGROQ,
		InferenceProviderOpenRouter,
		InferenceProviderCerebras,
		InferenceProviderVenice,
		InferenceProviderChutes,
		InferenceProviderHuggingFace,
		InferenceAIHubMix,
		InferenceKimiCoding,
		InferenceProviderCopilot,
		InferenceProviderVercel,
		InferenceProviderMiniMax,
	}
}

// KnownProviderTypes 返回所有已知的推理提供商类型
func KnownProviderTypes() []Type {
	return []Type{
		TypeOpenAI,
		TypeOpenAICompat,
		TypeOpenRouter,
		TypeVercel,
		TypeAnthropic,
		TypeGoogle,
		TypeAzure,
		TypeBedrock,
		TypeVertexAI,
	}
}
