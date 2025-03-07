package models

import (
	"slices"

	"github.com/openai/openai-go"
)

const (
	OpenAIModelGPT4oMini openai.ChatModel = openai.ChatModelGPT4oMini
	OpenAIModelGPT4o     openai.ChatModel = openai.ChatModelGPT4o
	OpenAIModelO3Mini    openai.ChatModel = openai.ChatModelO3Mini
)

func IsSupportedModel(model openai.ChatModel) bool {
	return slices.Contains(OpenAISupportedModels, model)
}

type Cost float64

// https://openai.com/api/pricing/
const (
	// GPT-4o Mini
	OpenAIModelGPT4oMiniInputCostPerToken       Cost = 2.50 * 1e-6
	OpenAIModelGPT4oMiniCachedInputCostPerToken Cost = 1.25 * 1e-6
	OpenAIModelGPT4oMiniOutputCostPerToken      Cost = 10.00 * 1e-6
	// GPT-4o
	OpenAIModelGPT4oInputCostPerToken       Cost = 0.150 * 1e-6
	OpenAIModelGPT4oCachedInputCostPerToken Cost = 0.075 * 1e-6
	OpenAIModelGPT4oOutputCostPerToken      Cost = 0.670 * 1e-6
	// O3 Mini
	OpenAIModelO3MiniInputCostPerToken       Cost = 1.10 * 1e-6
	OpenAIModelO3MiniCachedInputCostPerToken Cost = 0.55 * 1e-6
	OpenAIModelO3MiniOutputCostPerToken      Cost = 4.40 * 1e-6
)

type CostPerToken struct {
	Input       Cost
	CachedInput Cost
	Output      Cost
}

var OpenAIModelCosts = map[openai.ChatModel]CostPerToken{
	OpenAIModelGPT4oMini: {
		Input:       OpenAIModelGPT4oMiniInputCostPerToken,
		CachedInput: OpenAIModelGPT4oMiniCachedInputCostPerToken,
		Output:      OpenAIModelGPT4oMiniOutputCostPerToken,
	},
	OpenAIModelGPT4o: {
		Input:       OpenAIModelGPT4oInputCostPerToken,
		CachedInput: OpenAIModelGPT4oCachedInputCostPerToken,
		Output:      OpenAIModelGPT4oOutputCostPerToken,
	},
	OpenAIModelO3Mini: {
		Input:       OpenAIModelO3MiniInputCostPerToken,
		CachedInput: OpenAIModelO3MiniCachedInputCostPerToken,
		Output:      OpenAIModelO3MiniOutputCostPerToken,
	},
}

var OpenAISupportedModels = []openai.ChatModel{
	OpenAIModelGPT4oMini,
	OpenAIModelGPT4o,
	OpenAIModelO3Mini,
}

func EstimateCost(model openai.ChatModel, usage openai.CompletionUsage) Cost {
	cost := OpenAIModelCosts[model]
	estimatedCost := float64(cost.Input)*float64(usage.PromptTokens) +
		float64(cost.CachedInput)*float64(usage.PromptTokensDetails.CachedTokens) +
		float64(cost.Output)*float64(usage.CompletionTokens)
	return Cost(estimatedCost)
}
