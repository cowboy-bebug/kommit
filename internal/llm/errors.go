package llm

import "fmt"

type APIKeyMissingError struct{}
type OpenAIRequestError struct{ Err error }
type JSONParseError struct{ Err error }

func (e APIKeyMissingError) Error() string {
	return "KOMMIT_API_KEY or OPENAI_API_KEY environment variable must be set"
}

func (e OpenAIRequestError) Error() string {
	return fmt.Sprintf("OpenAI request failed: %v", e.Err)
}

func (e JSONParseError) Error() string {
	return fmt.Sprintf("JSON unmarshal failed: %v", e.Err)
}
