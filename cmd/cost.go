package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cowboy-bebug/kommit/internal/ui"
	"github.com/cowboy-bebug/kommit/internal/utils"
	"github.com/spf13/cobra"
)

var costCmd = &cobra.Command{
	Use:   "cost",
	Short: "💰 Show the cost of your commitment therapy",
	Long: `💰 Kommit Cost Analysis - Because emotional growth has a price tag!

This command reveals how much your repository's therapy sessions have cost you so far.
Think of it as reviewing your therapy bills - a necessary part of the healing process.

Just as real therapy isn't free, AI-powered commit messages come with a cost.
Still, it's cheaper than the emotional damage of trying to decipher "fixed stuff"
six months from now when you're debugging at 3 AM.

Your wallet might need its own support group, but your future self will thank you
for the investment in clear, meaningful commit history.`,
	Run: runCost,
}

func runCost(cmd *cobra.Command, args []string) {
	costs, err := utils.GetCosts()
	if err != nil {
		if errors.Is(err, utils.CostFileNotFoundError{}) {
			fmt.Println("😰 Financial abandonment detected: It hasn't committed any expenses yet.")
			fmt.Println("(Have you run `git kommit` yet?)")
			os.Exit(0)
		}
		fmt.Println("😰 Financial abandonment detected: Failed to retrieve your expenses.")
		os.Exit(1)
	}

	if err = ui.CostTableModel(costs); err != nil {
		ui.HandleQuitError(err)
		fmt.Println("😰 Financial abandonment detected: Failed to display your expenses.")
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(costCmd)
}
