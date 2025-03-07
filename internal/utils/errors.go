package utils

import (
	"fmt"

	"github.com/openai/openai-go"
)

type UnsupportedModelError struct{ Model openai.ChatModel }

func (e UnsupportedModelError) Error() string {
	return fmt.Sprintf("Unsupported model: %s", e.Model)
}

func (e UnsupportedModelError) Is(target error) bool {
	_, ok := target.(UnsupportedModelError)
	return ok
}
