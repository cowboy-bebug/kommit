package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/cowboy-bebug/kommit/internal/llm"
	"github.com/cowboy-bebug/kommit/internal/ui"
	"github.com/cowboy-bebug/kommit/internal/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const commitMessageSignature = "[Generated by Kommit]"

var rootCmd = &cobra.Command{
	Use:   "kommit",
	Short: "😌 Git therapy for your commitment issues",
	Long: `😌 Kommit - AI-powered therapy for repositories with commitment issues

"So, tell me about these changes you've made..."

Tired of staring at the blank commit message prompt? Let Kommit analyze your
changes and craft meaningful, conventional commit messages that tell the real
story of your code.

This command helps your staged changes express themselves through
AI-generated commit messages that follow conventional commit formats.
It's like couples therapy between you and your future self - improving
communication now to prevent confusion and frustration later.

No more commitment issues. No more vague messages. Just healthy,
expressive documentation of your development journey.`,
	Run: runCommit,
}

func runCommit(cmd *cobra.Command, args []string) {
	// Check if there are staged changes
	diff, err := utils.ExecGit("diff", "--cached")
	if err != nil || diff == "" {
		fmt.Println("😰 Commitment issues detected: You're not ready to commit... anything.")
		fmt.Println("(Stage some changes first!)")
		os.Exit(1)
	}

	// Load config to get available scopes
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Println("😰 Commitment issues detected: You haven't booked your first therapy session!")
		fmt.Println("(Run 'git kommit init' to get on the calendar.)")
		if Verbose {
			log.Printf("Error loading config: %v", err)
		}
		os.Exit(1)
	}

	context := "I'm using the following conventional commit types:\n"
	context += fmt.Sprintf("- types: %s\n", config.Commit.Types)
	context += "Optionally use the following scopes only if the changes are related to the scopes:\n"
	context += fmt.Sprintf("- scopes: %s\n", config.Commit.Scopes)

	s := utils.Spinner("🧐 Helping your code express its feelings to future developers...")
	s.Start()
	commitMessage, err := llm.GenerateCommitMessage(config, diff)
	s.Stop()
	if err != nil {
		fmt.Println("😰 Commitment issues detected: Your code is experiencing emotional resistance!")
		if errors.Is(err, &llm.APIKeyMissingError{}) {
			fmt.Println("\nHave you set up your OpenAI API key? Try one of these:")
			fmt.Println("  export OPENAI_API_KEY=\"sk-...\"")
			fmt.Println("  export KOMMIT_API_KEY=\"sk-...\"    # For a dedicated key")
		}
		if Verbose {
			log.Printf("Error generating commit message: %v", err)
		}
		os.Exit(1)
	}
	commitMessage += fmt.Sprintf("\n\n%s", commitMessageSignature)

	fmt.Println("💭 Your therapist's recommendation:")
	fmt.Println("```text")
	color.New(color.FgGreen, color.Bold).Println(commitMessage)
	fmt.Println("```")

	option, err := ui.SelectCommit()
	ui.HandleQuitError(err)
	if err != nil {
		log.Printf("Error confirming commit: %v", err)
		os.Exit(1)
	}

	switch option {
	case ui.CommitOptionProceed:
		fmt.Println("🧐 Preparing for your code's commitment ceremony...")

		tempFile, err := os.CreateTemp("", ".kommit-msg-*.txt")
		if err != nil {
			fmt.Println("😰 Commitment issues detected: Refusing to prepare temporary paperwork!")
			if Verbose {
				log.Printf("Error creating temp file: %v", err)
			}
			os.Exit(1)
		}
		tempFilePath := tempFile.Name()
		defer os.Remove(tempFilePath) // Clean up the temp file when done

		if _, err := tempFile.WriteString(commitMessage); err != nil {
			tempFile.Close()
			fmt.Println("😰 Commitment issues detected: Refusing to fill the temporary paperwork!")
			if Verbose {
				log.Printf("Error writing to temp file: %v", err)
			}
			os.Exit(1)
		}
		tempFile.Close()

		cmd := exec.Command("git", "commit", "-q", "-F", tempFilePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println("😰 Commitment issues detected: Refusing to commit!")
			if Verbose {
				log.Printf("Error committing: %v", err)
			}
			os.Exit(1)
		}
		fmt.Println("🧐 Successfully committed! Your relationship with the repo has deepened!")
	case ui.CommitOptionEdit:
		fmt.Println("🧐 Starting your self-guided therapy session...")

		tempFile, err := os.CreateTemp("", ".kommit-msg-*.txt")
		if err != nil {
			fmt.Println("😰 Commitment issues detected: Couldn't prepare your self-therapy materials!")
			if Verbose {
				log.Printf("Error creating temp file: %v", err)
			}
			os.Exit(1)
		}
		tempFilePath := tempFile.Name()
		defer os.Remove(tempFilePath)

		// Write the suggested message to the file as a starting point
		if _, err := tempFile.WriteString(commitMessage); err != nil {
			tempFile.Close()
			fmt.Println("😰 Commitment issues detected: Couldn't write your therapy starting notes!")
			if Verbose {
				log.Printf("Error writing to temp file: %v", err)
			}
			os.Exit(1)
		}
		tempFile.Close()

		// Set the GIT_EDITOR environment variable to open the editor with the file
		cmd := exec.Command("git", "commit", "--template", tempFilePath, "--allow-empty-message")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Println("📝 Opening your personal therapy journal (editor)...")

		err = cmd.Run()
		if err != nil {
			fmt.Println("😰 Commitment issues detected: Your self-therapy session was interrupted!")
			if Verbose {
				log.Printf("Error during commit: %v", err)
			}
			os.Exit(1)
		}

		fmt.Println("🎓 Self-therapy complete! You've committed to your own path of growth.")
	case ui.CommitOptionRerun:
		runCommit(cmd, args)
	case ui.CommitOptionExit:
		fmt.Println("🧐 You're on your own path now. Call if your commitment issues return!")
		os.Exit(0)
	}
}

var Verbose bool
var Debug bool

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	rootCmd.PersistentFlags().BoolP("help", "h", false,
		"Schedule an emergency therapy session (show help)")

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false,
		"Hear all the relationship details your repo normally keeps private")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
