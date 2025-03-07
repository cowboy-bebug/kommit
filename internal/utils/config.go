package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cowboy-bebug/kommit/internal/models"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const configFilename = ".kommitrc.yaml"

func GetConfigPath() (string, error) {
	output, err := ExecGit("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func GetConfigFilePath() (string, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(configPath, configFilename), nil
}

func unmarshalConfig(v *viper.Viper) (*Config, error) {
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	return &config, nil
}

type LLMConfig struct {
	Model string `mapstructure:"model"`
}

type CommitConfig struct {
	Types  []string `mapstructure:"types"`
	Scopes []string `mapstructure:"scopes"`
}

type Config struct {
	LLM    LLMConfig    `mapstructure:"llm"`
	Commit CommitConfig `mapstructure:"commit"`
}

func LoadConfig() (*Config, error) {
	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetConfigFile(configFilePath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	config, err := unmarshalConfig(v)
	if err != nil {
		return nil, err
	}

	if !models.IsSupportedModel(config.LLM.Model) {
		return nil, UnsupportedModelError{Model: config.LLM.Model}
	}

	return config, nil
}

func GetDefaultConfig() (*Config, error) {
	v := viper.New()
	v.SetDefault("llm", map[string]any{
		"model": models.OpenAIModelGPT4oMini,
	})
	v.SetDefault("commit", map[string]any{
		"types": []string{
			"build",
			"chore",
			"ci",
			"docs",
			"feat",
			"fix",
			"perf",
			"refactor",
			"revert",
			"style",
			"test",
		},
		"scopes": []string{},
	})
	return unmarshalConfig(v)
}

func createSequenceNode(items []string) *yaml.Node {
	node := &yaml.Node{
		Kind: yaml.SequenceNode,
	}

	for _, item := range items {
		node.Content = append(node.Content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: item,
		})
	}

	return node
}

func WriteConfig(config *Config) error {
	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	// Create a file to write to
	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}
	defer file.Close()

	// Create a YAML encoder with the desired indentation
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2) // Set indentation to 2 spaces

	// Create a yaml.Node to control the order of fields
	root := &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					// `llm` key
					{Kind: yaml.ScalarNode, Value: "llm"},
					{
						Kind: yaml.MappingNode,
						Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Value: "model"},
							{Kind: yaml.ScalarNode, Value: config.LLM.Model},
						},
					},

					// `commit` key
					{Kind: yaml.ScalarNode, Value: "commit"},
					{
						Kind: yaml.MappingNode,
						Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Value: "types"},
							createSequenceNode(config.Commit.Types),
							{Kind: yaml.ScalarNode, Value: "scopes"},
							createSequenceNode(config.Commit.Scopes),
						},
					},
				},
			},
		},
	}

	// Encode the node structure
	if err := encoder.Encode(root); err != nil {
		return fmt.Errorf("error encoding config to YAML: %w", err)
	}

	return nil
}

func PrintConfigFile() error {
	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("💊 Prescription for better commit messages written to: %s\n", configFilePath)
	fmt.Println("---------------------------")
	fmt.Println(strings.TrimRight(string(data), "\n"))
	fmt.Println("---------------------------")

	return nil
}
