package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string
	Commit  string
	Date    string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "🏷️ Show the therapist's credentials",
	Long: `🏷️ Kommitment Therapist Credentials - Because even AI therapists need qualifications!

This command reveals the version information about your Kommitment therapist,
including their training date and certification level.

Just as you wouldn't trust your emotional well-being to an unqualified therapist,
your repository deserves a properly versioned commit message counselor.

Your therapist has been helping repositories express themselves meaningfully
since their certification date, and continues to grow with each update.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Kommitment %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", Date)
		fmt.Printf("\n🧐 Therapist v%s: Licensed to treat repositories with severe commitment issues since %s.\n",
			Version, Date[:10])
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
