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
		return nil, &APIKeyMissingError{}
	}

	return openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithRequestTimeout(timeout),
	), nil
}

func chat(model, prompt string) (string, error) {
	client, err := newClient()
	if err != nil {
		return "", err
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
		return "", &OpenAIRequestError{Err: err}
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
		return empty, err
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
		return empty, &OpenAIRequestError{Err: err}
	}

	content := resp.Choices[0].Message.Content
	var result T
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		var empty T
		return empty, &JSONParseError{Err: err}
	}

	return result, nil
}

func GenerateCommitMessage(model, context, diff string) (string, error) {
	// main prompt
	prompt := "Generate a single commit message in the conventional commit message format"

	// constraints
	prompt += ", without:\n"
	prompt += "- wrapping in a code block\n"
	prompt += "- using build as scope\n"
	prompt += "- suggesting `feat` for build scripts\n"
	prompt += "- suggesting comments or remarks\n"

	// type
	prompt += "\nUsing conventional commit types:\n"
	prompt += "- in lowercase\n"
	prompt += "- Use `build` for build system, scripts or settings, such as Makefile, Dockerfile, etc.\n"
	prompt += "- Use `docs` for changes to the documentation, such as README, CHANGELOG, etc.\n"
	prompt += "- Do not use `docs` for changes to scripts or code\n"

	// scope
	prompt += "\nUsing conventional commit scopes:\n"
	prompt += "- use the module or package name\n"
	prompt += "- leave empty if the changes are not related to a specific module or package\n"
	prompt += "- leave empty if the changes are across multiple modules or packages\n"

	// subject
	prompt += "\nUsing conventional commit message subject:\n"
	prompt += "- using the imperative mood (present tense)\n"

	// message body
	prompt += "\nUsing conventional commit message body:\n"
	prompt += "- as bullet points for the changes\n"
	prompt += "- using the imperative mood (present tense)\n"
	prompt += "- using titlecase for the first letter of the message body\n"
	prompt += "- breaking lines at 72 characters\n"
	prompt += "- Generate a message body only if the changes are significant\n"

	// context and diff
	prompt += "\nBased on the following context and diff:\n"
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
