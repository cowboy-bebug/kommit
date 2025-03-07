package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cowboy-bebug/kommit/internal/models"
	"github.com/cowboy-bebug/kommit/internal/utils"
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAI client configuration
const (
	timeout          = 10 * time.Second
	temperature      = 0.0
	topP             = 1.0
	presencePenalty  = 0.0
	frequencyPenalty = 0.0
)

// System prompts
const (
	kommitSystemPrompt = "You are an AI that generates Conventional Git commit messages."
	jsonResponsePrompt = "Return your response as a valid JSON object."
)

// User prompts
const (
	promptMain = "Generate a single commit message following the **Conventional Commit** format, adhering to these rules:\n"

	promptGeneralRules = `
## **General Rules**
- **Do not**:
  - Wrap the message in a code block or triple backticks.
  - Use ` + "`" + "build" + "`" + ` as a scope.
  - Suggest ` + "`" + "feat" + "`" + ` for build scripts.
  - Include comments or remarks.

- **Do**:
  - Try your best to guess what the git diff is about.
  - Wrap lines at **72 characters**.
`

	promptCommitTypeGuidelines = `
## **Commit Type Guidelines**
- Use **lowercase** commit types:
  - ` + "`" + "build" + "`" + `: For build systems, scripts, or settings (e.g., Makefile, Dockerfile).
  - ` + "`" + "docs" + "`" + `: For documentation changes (e.g., README, CHANGELOG), **but not** script or code changes.
`

	promptScopeRules = `
## **Scope Rules**
- Use the **module or package name** as the scope.
- **Leave the scope empty** if:
  - The changes are **not** tied to a specific module or package.
  - The changes span **multiple modules, packages, files or scopes**.
`

	promptMessageFormatting = `
## **Message Formatting**
- **Subject**:
  - Use **imperative mood** (present tense).
- **Body _(only if changes are significant)_:
  - Use **bullet points**.
  - Use **imperative mood** (present tense).
  - Capitalize the **first letter** of each bullet point.
  - Wrap lines at **72 characters**.
`

	kommitBaseUserPrompt = promptMain + promptGeneralRules + promptCommitTypeGuidelines + promptScopeRules + promptMessageFormatting
)

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

type ChatResult[T any] struct {
	Message T
	Cost    models.Cost
}

func chat(model, prompt string) (ChatResult[string], error) {
	client, err := newClient()
	if err != nil {
		return ChatResult[string]{}, err
	}

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.F(model),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(kommitSystemPrompt),
			openai.UserMessage(prompt),
		}),
		Temperature:      openai.Float(temperature),
		TopP:             openai.Float(topP),
		PresencePenalty:  openai.Float(presencePenalty),
		FrequencyPenalty: openai.Float(frequencyPenalty),
	})
	if err != nil {
		return ChatResult[string]{}, &OpenAIRequestError{Err: err}
	}

	return ChatResult[string]{
		Message: resp.Choices[0].Message.Content,
		Cost:    models.EstimateCost(model, resp.Usage),
	}, nil
}

func wrapInCSVCodeBlock(x []string) string {
	l := make([]string, len(x))
	for i, s := range x {
		l[i] = fmt.Sprintf("`%s`", s)
	}
	return fmt.Sprintf("  - %s\n", strings.Join(l, ", "))
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

func chatStructured[T any](model, prompt string, schema openai.ResponseFormatJSONSchemaJSONSchemaParam) (ChatResult[T], error) {
	client, err := newClient()
	if err != nil {
		return ChatResult[T]{}, err
	}

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model:            openai.F(model),
		Temperature:      openai.Float(temperature),
		TopP:             openai.Float(topP),
		PresencePenalty:  openai.Float(presencePenalty),
		FrequencyPenalty: openai.Float(frequencyPenalty),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(kommitSystemPrompt + jsonResponsePrompt),
			openai.UserMessage(prompt),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schema),
			}),
	})
	if err != nil {
		return ChatResult[T]{}, &OpenAIRequestError{Err: err}
	}

	content := resp.Choices[0].Message.Content
	var result T
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return ChatResult[T]{}, &JSONParseError{Err: err}
	}

	return ChatResult[T]{
		Message: result,
		Cost:    models.EstimateCost(model, resp.Usage),
	}, nil
}

func GenerateCommitMessage(config *utils.Config, diff string) (ChatResult[string], error) {
	prompt := kommitBaseUserPrompt

	// context: commit types
	prompt += "\n## Context:\n"
	prompt += "- **Allowed commit types**:\n"
	prompt += wrapInCSVCodeBlock(config.Commit.Types)

	// context: commit scopes
	prompt += "- **Allowed scopes _(only if changes are limited to a single scope)_:\n"
	prompt += wrapInCSVCodeBlock(config.Commit.Scopes)
	prompt += "  - **Note:** If the changes span multiple scopes, do not use a scope in the commit message.\n"

	// diff
	prompt += "\n## Git Diff:\n"
	prompt += "**Based on the following diff**:\n"
	prompt += "```diff\n"
	prompt += diff + "\n"
	prompt += "```\n"

	return chat(config.LLM.Model, prompt)
}

type Scopes struct {
	Scopes []string `json:"scopes"`
}

var StructuredScopesSchema = GenerateSchema[Scopes]()

func GenerateScopesFromFilenames(model string, filenames, existingScopes []string) (ChatResult[Scopes], error) {
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
		return ChatResult[Scopes]{}, err
	}

	return result, nil
}
