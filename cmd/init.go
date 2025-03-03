package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cowboy-bebug/kommit/internal/llm"
	"github.com/cowboy-bebug/kommit/internal/utils"
	"github.com/manifoldco/promptui"
	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "😌 Your repo's therapy session begins here",
	Long: `😌 Initialize Kommit - because your repo has commitment issues!

This command creates your .kommitrc.yaml file with the emotional intelligence
your git history desperately needs. It'll analyze your project structure and
suggest meaningful scopes so your commits can finally express themselves properly.

Remember: Good commit messages are like good apologies - specific, sincere,
and they don't include the phrase "various changes". Your future self will
thank you for the therapy.`,
	Run: runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	// Load config, return if it exists
	config, err := utils.LoadConfig()
	if config != nil {
		fmt.Println("🥹 Your repo is already in therapy! Treatment plan exists.")
		fmt.Println("🥰 Run `kommit commit` to continue the healing process!")
		os.Exit(0)
	}

	// Get default config, if it doesn't exist
	if err != nil {
		config, err = utils.GetDefaultConfig()
		if err != nil {
			fmt.Println("😰 Therapy session interrupted: Failed to retrieve your treatment plan.")
			fmt.Println()
			if Verbose {
				log.Printf("Error getting default config: %v", err)
			}
			os.Exit(1)
		}
	}

	// Select model
	prompt := promptui.Select{
		Label: "Choose your therapist's qualifications:",
		Items: []string{
			openai.ChatModelGPT4oMini,
			openai.ChatModelGPT4o,
			openai.ChatModelO3Mini,
		},
	}
	_, model, err := prompt.Run()
	if err != nil {
		fmt.Println("😰 Therapy session interrupted: Failed to select your therapist.")
		fmt.Println()
		if Verbose {
			log.Printf("Error selecting model: %v", err)
		}
		os.Exit(1)
	}
	config.LLM.Model = model

	// Get scopes from commit history
	existingScopes, err := utils.GetScopesFromHistory()
	if err != nil {
		if Verbose {
			log.Printf("Error getting scopes from git history: %v", err)
		}
	}

	// Get files from directory
	filenames, err := utils.GetFilesFromDirectory(5)
	if err != nil {
		fmt.Println("😰 Therapy session interrupted: Failed to retrieve your relationship history.")
		fmt.Println()
		if Verbose {
			log.Printf("Error getting scopes from directory: %v", err)
		}
		os.Exit(1)
	}

	// Generate scopes from directory
	scopes, err := llm.GenerateScopesFromFilenames(model, filenames, existingScopes)
	if err != nil {
		fmt.Println("😰 Therapy session interrupted: Failed to establish your treatment plan.")
		fmt.Println()
		if Verbose {
			log.Printf("Error generating scopes from directory: %v", err)
		}
		os.Exit(1)
	}

	// Write config
	config.Commit.Scopes = scopes
	err = utils.WriteConfig(config)
	if err != nil {
		fmt.Println("😰 Therapy session interrupted: Failed to write your treatment plan.")
		fmt.Println()
		if Verbose {
			log.Printf("Error writing config: %v", err)
		}
		os.Exit(1)
	}

	fmt.Println("🥹 Your repo is in therapy! Treatment plan filled successfully.")
	fmt.Println("🥰 Run `kommit commit` to continue the healing process!")
	fmt.Println()
	utils.PrintConfigFile()
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(initCommand)
}
