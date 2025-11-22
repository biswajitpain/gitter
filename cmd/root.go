package cmd

import (
	"fmt"
	"io"
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
		ExecuteRootCommand(cmd, args, os.Stdout, os.Stderr)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}

// ExecuteRootCommand is a helper function to allow testing rootCmd.Run with custom stdout/stderr
func ExecuteRootCommand(cmd *cobra.Command, args []string, stdout, stderr io.Writer) int {
	// Check if 'git' is available.
	if _, err := execCommand("git", "version").Output(); err != nil { // Use execCommand
		fmt.Fprintf(stderr, "Error: 'git' command not found. Please ensure Git is installed and in your PATH.\n")
		return 1
	}

	// If no args are passed, show our own help.
	if len(args) == 0 {
		cmd.SetOut(stdout) // Set Cobra's output for help message
		cmd.SetErr(stderr) // Set Cobra's error output for help message
		cmd.Help()
		return 0
	}

	// If we reach here, it means Cobra didn't find a matching subcommand.
	// So, we assume it's a git command and pass all arguments to git.
	gitCmd := execCommand("git", args...)
	gitCmd.Stdin = os.Stdin
	gitCmd.Stdout = stdout
	gitCmd.Stderr = stderr

	if err := gitCmd.Run(); err != nil {
		// Exit with the same code as the git command
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		} else {
			return 1
		}
	}
	return 0
}

func init() {
	// Disable Cobra's default error and usage printing for unknown commands.
	// This allows the Run function of rootCmd to handle the passthrough to git.
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
}
