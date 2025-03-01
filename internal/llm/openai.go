package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const timeout = 10 * time.Second
const temperature = 0.0
const topP = 1.0
const presencePenalty = 0.0
const frequencyPenalty = 0.0

const kommitPrompt = "You are an AI that generates Conventional Git commit messages."
const jsonResponsePrompt = "Return your response as a valid JSON object."

func newClient() (*openai.Client, error) {
	// KOMMIT_API_KEY takes precedence
	apiKey := os.Getenv("KOMMIT_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	// If both are missing, return an error
	if apiKey == "" {
		return nil, fmt.Errorf("KOMMIT_API_KEY or OPENAI_API_KEY environment variable must be set")
	}

	return openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithRequestTimeout(timeout),
	), nil
}

func chat(model, prompt string) (string, error) {
	client, err := newClient()
	if err != nil {
		return "", fmt.Errorf("error creating OpenAI client: %v", err)
	}

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.F(model),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(kommitPrompt),
			openai.UserMessage(prompt),
		}),
		Temperature:      openai.Float(temperature),
		TopP:             openai.Float(topP),
		PresencePenalty:  openai.Float(presencePenalty),
		FrequencyPenalty: openai.Float(frequencyPenalty),
	})
	if err != nil {
		return "", fmt.Errorf("OpenAI request failed: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func GenerateSchema[T any]() any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func chatStructured[T any](model, prompt string, schema openai.ResponseFormatJSONSchemaJSONSchemaParam) (T, error) {
	client, err := newClient()
	if err != nil {
		var empty T
		return empty, fmt.Errorf("error creating OpenAI client: %v", err)
	}

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model:            openai.F(model),
		Temperature:      openai.Float(temperature),
		TopP:             openai.Float(topP),
		PresencePenalty:  openai.Float(presencePenalty),
		FrequencyPenalty: openai.Float(frequencyPenalty),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(kommitPrompt + jsonResponsePrompt),
			openai.UserMessage(prompt),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schema),
			}),
	})
	if err != nil {
		var empty T
		return empty, fmt.Errorf("OpenAI request failed: %v", err)
	}

	content := resp.Choices[0].Message.Content
	var result T
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		var empty T
		return empty, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return result, nil
}

func GenerateCommitMessage(model, context, diff string) (string, error) {
	// constraints
	prompt := "Without wrapping in a code block and"
	prompt += " without suggesting comments or remarks"

	// main prompt
	prompt += ", can you generate a conventional commit message"

	// message body
	prompt += ", with message body as bullet points for the changes"
	prompt += ", and breaking lines at 72 characters"

	// context and diff
	prompt += ", based on the following context and diff:\n"
	prompt += fmt.Sprintf("Here is context: %s\n", context)
	prompt += fmt.Sprintf("Here is the Git diff:\n%s\n", diff)

	return chat(model, prompt)
}

type Scopes struct {
	Scopes []string `json:"scopes"`
}

var StructuredScopesSchema = GenerateSchema[Scopes]()

func GenerateScopesFromFilenames(model string, filenames, existingScopes []string) ([]string, error) {
	prompt := "Based on the following project structure, guess module or package names used in this project:\n"
	prompt += strings.Join(filenames, "\n")

	prompt += "Here are some existing scopes:\n"
	prompt += strings.Join(existingScopes, "\n")

	prompt += "\n\n"
	prompt += "- Do not suggest nested names\n"
	prompt += "- Do not suggest names with \"/\""
	prompt += "- Do not suggest docs as a scope"

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("names"),
		Description: openai.F("A list of module or package names."),
		Schema:      openai.F(StructuredScopesSchema),
		Strict:      openai.Bool(true),
	}

	result, err := chatStructured[Scopes](model, prompt, schemaParam)
	if err != nil {
		return nil, err
	}

	return result.Scopes, nil
}
