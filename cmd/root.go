package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitter",
	Short: "gitter is a smart git wrapper",
	Long: `gitter is a command-line wrapper for Git that enhances your commit workflow 
with an intelligent "commit review" (cr) command and LLM-powered message generation.`,
	// This is the core of the git passthrough logic.
	// We disable flag parsing for the root command and check if the
	// command is one of our own. If not, we pass it to git.
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if 'git' is available.
		if _, err := exec.LookPath("git"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: 'git' command not found. Please ensure Git is installed and in your PATH.\n")
			os.Exit(1)
		}

		// If no args are passed, show our own help.
		if len(args) == 0 {
			cmd.Help()
			return
		}

		// Check if the command is one of ours. If so, Cobra has already handled it.
		// This Run function only executes for commands that are NOT registered with Cobra.
		// So, we can assume it's a git command.
		gitCmd := exec.Command("git", args...)
		gitCmd.Stdin = os.Stdin
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			// Exit with the same code as the git command
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			} else {
				os.Exit(1)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}
