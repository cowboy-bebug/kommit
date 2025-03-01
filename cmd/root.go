package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kommit",
	Short: "😌 Git therapy for your commitment issues",
	Long: `😌 Kommitment - AI-powered therapy for repositories with commitment issues

Tired of staring at the blank commit message prompt? Let Kommit analyze your
changes and craft meaningful, conventional commit messages that tell the real
story of your code.

Because every good relationship needs clear communication - even the one
between you and your future self reading these commits.`,
}

var Verbose bool

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
