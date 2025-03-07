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

var OpenAISupportedModels = []openai.ChatModel{
	OpenAIModelGPT4oMini,
	OpenAIModelGPT4o,
	OpenAIModelO3Mini,
}

func IsSupportedModel(model openai.ChatModel) bool {
	return slices.Contains(OpenAISupportedModels, model)
}
