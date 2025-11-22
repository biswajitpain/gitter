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
	// Disable Cobra's default "unknown command" error and let our Run function handle it.
	// This allows us to pass unknown commands directly to git.
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

		// If we reach here, it means Cobra didn't find a matching subcommand.
		// So, we assume it's a git command and pass all arguments to git.
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
	// Disable Cobra's default error and usage printing for unknown commands.
	// This allows the Run function of rootCmd to handle the passthrough to git.
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
}
