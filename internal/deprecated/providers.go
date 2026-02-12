// Package deprecated 用于提供旧版本的AI提供商配置
package deprecated

import (
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed configs/openai.json
var openAIConfig []byte // OpenAI配置文件内容

//go:embed configs/anthropic.json
var anthropicConfig []byte // Anthropic配置文件内容

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

// ProviderFunc 是一个返回Provider的函数类型
type ProviderFunc func() Provider

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
	groqProvider,
	openRouterProvider,
	cerebrasProvider,
	veniceProvider,
	chutesProvider,
	deepSeekProvider,
	huggingFaceProvider,
	aiHubMixProvider,
}

// GetAll 返回所有已注册的AI提供商配置
func GetAll() []Provider {
	providers := make([]Provider, 0, len(providerRegistry))
	for _, providerFunc := range providerRegistry {
		providers = append(providers, providerFunc())
	}
	return providers
}

// loadProviderFromConfig 从JSON配置数据中加载提供商配置
func loadProviderFromConfig(configData []byte) Provider {
	var p Provider
	if err := json.Unmarshal(configData, &p); err != nil {
		log.Printf("加载提供商配置出错: %v", err)
		return Provider{}
	}
	return p
}

// openAIProvider 返回OpenAI提供商配置
func openAIProvider() Provider {
	return loadProviderFromConfig(openAIConfig)
}

// anthropicProvider 返回Anthropic提供商配置
func anthropicProvider() Provider {
	return loadProviderFromConfig(anthropicConfig)
}

// geminiProvider 返回Gemini提供商配置
func geminiProvider() Provider {
	return loadProviderFromConfig(geminiConfig)
}

// azureProvider 返回Azure提供商配置
func azureProvider() Provider {
	return loadProviderFromConfig(azureConfig)
}

// bedrockProvider 返回Bedrock提供商配置
func bedrockProvider() Provider {
	return loadProviderFromConfig(bedrockConfig)
}

// vertexAIProvider 返回VertexAI提供商配置
func vertexAIProvider() Provider {
	return loadProviderFromConfig(vertexAIConfig)
}

// xAIProvider 返回xAI提供商配置
func xAIProvider() Provider {
	return loadProviderFromConfig(xAIConfig)
}

// zAIProvider 返回zAI提供商配置
func zAIProvider() Provider {
	return loadProviderFromConfig(zAIConfig)
}

// openRouterProvider 返回OpenRouter提供商配置
func openRouterProvider() Provider {
	return loadProviderFromConfig(openRouterConfig)
}

// groqProvider 返回Groq提供商配置
func groqProvider() Provider {
	return loadProviderFromConfig(groqConfig)
}

// cerebrasProvider 返回Cerebras提供商配置
func cerebrasProvider() Provider {
	return loadProviderFromConfig(cerebrasConfig)
}

// veniceProvider 返回Venice提供商配置
func veniceProvider() Provider {
	return loadProviderFromConfig(veniceConfig)
}

// chutesProvider 返回Chutes提供商配置
func chutesProvider() Provider {
	return loadProviderFromConfig(chutesConfig)
}

// deepSeekProvider 返回DeepSeek提供商配置
func deepSeekProvider() Provider {
	return loadProviderFromConfig(deepSeekConfig)
}

// huggingFaceProvider 返回HuggingFace提供商配置
func huggingFaceProvider() Provider {
	return loadProviderFromConfig(huggingFaceConfig)
}

// aiHubMixProvider 返回AIHubMix提供商配置
func aiHubMixProvider() Provider {
	return loadProviderFromConfig(aiHubMixConfig)
}
