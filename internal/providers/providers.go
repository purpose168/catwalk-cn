// Package providers 提供推理提供商的注册管理功能
package providers

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/purpose168/catwalk-cn/pkg/catwalk"
)

//go:embed configs/openai.json
var openAIConfig []byte // OpenAI配置文件内容

//go:embed configs/anthropic.json
var anthropicConfig []byte // Anthropic配置文件内容

//go:embed configs/synthetic.json
var syntheticConfig []byte // Synthetic配置文件内容

//go:embed configs/gemini.json
var geminiConfig []byte // Gemini配置文件内容

//go:embed configs/openrouter.json
var openRouterConfig []byte // OpenRouter配置文件内容

//go:embed configs/azure.json
var azureConfig []byte // Azure配置文件内容

//go:embed configs/vertexai.json
var vertexAIConfig []byte // VertexAI配置文件内容

//go:embed configs/xai.json
var xAIConfig []byte // xAI配置文件内容

//go:embed configs/zai.json
var zAIConfig []byte // zAI配置文件内容

//go:embed configs/bedrock.json
var bedrockConfig []byte // Bedrock配置文件内容

//go:embed configs/groq.json
var groqConfig []byte // Groq配置文件内容

//go:embed configs/cerebras.json
var cerebrasConfig []byte // Cerebras配置文件内容

//go:embed configs/venice.json
var veniceConfig []byte // Venice配置文件内容

//go:embed configs/chutes.json
var chutesConfig []byte // Chutes配置文件内容

//go:embed configs/deepseek.json
var deepSeekConfig []byte // DeepSeek配置文件内容

//go:embed configs/huggingface.json
var huggingFaceConfig []byte // HuggingFace配置文件内容

//go:embed configs/aihubmix.json
var aiHubMixConfig []byte // AIHubMix配置文件内容

//go:embed configs/kimi.json
var kimiCodingConfig []byte // Kimi Coding配置文件内容

//go:embed configs/copilot.json
var copilotConfig []byte // Copilot配置文件内容

//go:embed configs/vercel.json
var vercelConfig []byte // Vercel配置文件内容

//go:embed configs/minimax.json
var miniMaxConfig []byte // MiniMax配置文件内容

// ProviderFunc 是一个返回Provider的函数类型
type ProviderFunc func() catwalk.Provider

// providerRegistry 存储所有注册的提供商构造函数
var providerRegistry = []ProviderFunc{
	anthropicProvider,
	openAIProvider,
	geminiProvider,
	azureProvider,
	bedrockProvider,
	vertexAIProvider,
	xAIProvider,
	zAIProvider,
	kimiCodingProvider,
	groqProvider,
	openRouterProvider,
	cerebrasProvider,
	veniceProvider,
	chutesProvider,
	deepSeekProvider,
	huggingFaceProvider,
	aiHubMixProvider,
	syntheticProvider,
	copilotProvider,
	vercelProvider,
	miniMaxProvider,
}

// GetAll 返回所有已注册的AI提供商配置
func GetAll() []catwalk.Provider {
	providers := make([]catwalk.Provider, 0, len(providerRegistry))
	for _, providerFunc := range providerRegistry {
		providers = append(providers, providerFunc())
	}
	return providers
}

// loadProviderFromConfig 从JSON配置数据中加载提供商配置
func loadProviderFromConfig(configData []byte) catwalk.Provider {
	var p catwalk.Provider
	if err := json.Unmarshal(configData, &p); err != nil {
		log.Printf("加载提供商配置出错: %v", err)
		return catwalk.Provider{}
	}
	return p
}

// openAIProvider 返回OpenAI提供商配置
func openAIProvider() catwalk.Provider {
	return loadProviderFromConfig(openAIConfig)
}

// syntheticProvider 返回Synthetic提供商配置
func syntheticProvider() catwalk.Provider {
	return loadProviderFromConfig(syntheticConfig)
}

// anthropicProvider 返回Anthropic提供商配置
func anthropicProvider() catwalk.Provider {
	return loadProviderFromConfig(anthropicConfig)
}

// geminiProvider 返回Gemini提供商配置
func geminiProvider() catwalk.Provider {
	return loadProviderFromConfig(geminiConfig)
}

// azureProvider 返回Azure提供商配置
func azureProvider() catwalk.Provider {
	return loadProviderFromConfig(azureConfig)
}

// bedrockProvider 返回Bedrock提供商配置
func bedrockProvider() catwalk.Provider {
	return loadProviderFromConfig(bedrockConfig)
}

// vertexAIProvider 返回VertexAI提供商配置
func vertexAIProvider() catwalk.Provider {
	return loadProviderFromConfig(vertexAIConfig)
}

// xAIProvider 返回xAI提供商配置
func xAIProvider() catwalk.Provider {
	return loadProviderFromConfig(xAIConfig)
}

// zAIProvider 返回zAI提供商配置
func zAIProvider() catwalk.Provider {
	return loadProviderFromConfig(zAIConfig)
}

// openRouterProvider 返回OpenRouter提供商配置
func openRouterProvider() catwalk.Provider {
	return loadProviderFromConfig(openRouterConfig)
}

// groqProvider 返回Groq提供商配置
func groqProvider() catwalk.Provider {
	return loadProviderFromConfig(groqConfig)
}

// cerebrasProvider 返回Cerebras提供商配置
func cerebrasProvider() catwalk.Provider {
	return loadProviderFromConfig(cerebrasConfig)
}

// veniceProvider 返回Venice提供商配置
func veniceProvider() catwalk.Provider {
	return loadProviderFromConfig(veniceConfig)
}

// chutesProvider 返回Chutes提供商配置
func chutesProvider() catwalk.Provider {
	return loadProviderFromConfig(chutesConfig)
}

// deepSeekProvider 返回DeepSeek提供商配置
func deepSeekProvider() catwalk.Provider {
	return loadProviderFromConfig(deepSeekConfig)
}

// huggingFaceProvider 返回HuggingFace提供商配置
func huggingFaceProvider() catwalk.Provider {
	return loadProviderFromConfig(huggingFaceConfig)
}

// aiHubMixProvider 返回AIHubMix提供商配置
func aiHubMixProvider() catwalk.Provider {
	return loadProviderFromConfig(aiHubMixConfig)
}

// kimiCodingProvider 返回Kimi Coding提供商配置
func kimiCodingProvider() catwalk.Provider {
	return loadProviderFromConfig(kimiCodingConfig)
}

// copilotProvider 返回Copilot提供商配置
func copilotProvider() catwalk.Provider {
	return loadProviderFromConfig(copilotConfig)
}

// vercelProvider 返回Vercel提供商配置
func vercelProvider() catwalk.Provider {
	return loadProviderFromConfig(vercelConfig)
}

// miniMaxProvider 返回MiniMax提供商配置
func miniMaxProvider() catwalk.Provider {
	return loadProviderFromConfig(miniMaxConfig)
}
