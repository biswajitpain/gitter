package cmd

import (
	"bufio"
	"fmt"
	"github.com/biswajitpain/gitter/internal/config"
	"github.com/biswajitpain/gitter/internal/llm"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// crCmd represents the cr command
var crCmd = &cobra.Command{
	Use:   "cr",
	Short: "Create a commit with an AI-generated message",
	Long: `The 'cr' command automates the commit process.
It stages files, generates a diff, and uses an LLM (if configured) 
to create a conventional commit message.`,
	Run: handleCrCommand,
}

func init() {
	rootCmd.AddCommand(crCmd)
}

// fileChangeStats holds the statistics for a single changed file.
type fileChangeStats struct {
	filePath     string
	linesChanged int
	charsChanged int
}

func handleCrCommand(cmd *cobra.Command, args []string) {
	// 1. Check if we are in a git repository.
	gitCheckCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := gitCheckCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Not a git repository. (%v)\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	// 2. Check for staged changes.
	stagedCheckCmd := exec.Command("git", "diff", "--cached", "--quiet")
	if stagedCheckCmd.Run() != nil { // Exits with 1 if no staged changes
		fmt.Print("No files are currently staged. Would you like to stage all changed files? (y/n): ")
		stageAllInput, _ := reader.ReadString('\n')
		stageAllInput = strings.TrimSpace(strings.ToLower(stageAllInput))

		if stageAllInput == "y" {
			fmt.Println("Staging all changed files...")
			if err := exec.Command("git", "add", ".").Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error staging changes: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Operation cancelled. No files were staged.")
			os.Exit(0)
		}
	} else {
		fmt.Println("Working on currently staged changes.")
	}

	// 3. Get the diff of staged changes.
	fmt.Println("Generating diff...")
	diffCmd := exec.Command("git", "diff", "--staged")
	diffOutputBytes, err := diffCmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting diff: %v\n", err)
		os.Exit(1)
	}
	diffOutput := string(diffOutputBytes)

	if strings.TrimSpace(diffOutput) == "" {
		fmt.Println("No changes to commit.")
		exec.Command("git", "reset").Run()
		os.Exit(0)
	}

	// 4. Parse the diff to get stats.
	stats := parseDiffStats(diffOutput)

	// 5. Ask the user for a commit message.
	fmt.Print("Please enter a commit message (or press Enter for a default):\n> ")
	userMessage, _ := reader.ReadString('\n')
	userMessage = strings.TrimSpace(userMessage)

	if userMessage == "" {
		userMessage = createDefaultCommitMessage()
		fmt.Printf("No commit message provided. Using default: \"%s\"\n", userMessage)
	}

	// 6. Generate a nice commit message.
	generatedMessage := generateCommitMessage(userMessage, diffOutput, stats)

	// 7. Ask for confirmation.
	fmt.Println("\n--- Generated Commit Message ---")
	fmt.Print(generatedMessage)
	fmt.Println("\n--------------------------------")
	fmt.Print("Confirm commit with this message? (y/n): ")
	confirmInput, _ := reader.ReadString('\n')
	confirmInput = strings.TrimSpace(strings.ToLower(confirmInput))

	if confirmInput == "y" {
		// 8. Commit.
		fmt.Println("Committing...")
		commitCmd := exec.Command("git", "commit", "-m", generatedMessage)
		if err := commitCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error committing: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Commit successful.")
	} else {
		fmt.Println("Commit cancelled. Changes are still staged.")
		fmt.Print("Would you like to unstage the changes? (y/n): ")
		unstageInput, _ := reader.ReadString('\n')
		unstageInput = strings.TrimSpace(strings.ToLower(unstageInput))
		if unstageInput == "y" {
			if err := exec.Command("git", "reset").Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error unstaging changes: %v\n", err)
			} else {
				fmt.Println("Changes have been unstaged.")
			}
		}
		os.Exit(0)
	}
}

func parseDiffStats(diffOutput string) []fileChangeStats {
	fileDiffs := strings.Split(diffOutput, "diff --git")
	var stats []fileChangeStats
	for _, fileDiff := range fileDiffs {
		if strings.TrimSpace(fileDiff) == "" {
			continue
		}

	
		lines := strings.Split(fileDiff, "\n")
		var currentFile string
		var linesChanged, charsChanged int

		for _, line := range lines {
			if strings.HasPrefix(line, "---") {
				currentFile = strings.TrimPrefix(line, "--- a/")
				break
			}
		}
		if currentFile == "" {
			continue
		}

		for _, line := range lines {
			if (strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "++")) ||
				(strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---")) {
				linesChanged++
				charsChanged += len(strings.TrimSpace(line[1:]))
			}
		}
		stats = append(stats, fileChangeStats{
			filePath:     currentFile,
			linesChanged: linesChanged,
			charsChanged: charsChanged,
		})
	}
	return stats
}

func createDefaultCommitMessage() string {
	repoPath, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not determine repo name for default commit message: %v\n", err)
		return "new git commit" // Fallback default
	}
	repoName := filepath.Base(strings.TrimSpace(string(repoPath)))

dateStr := time.Now().Format("2006-01-02")
	return fmt.Sprintf("new git commit on %s %s", repoName, dateStr)
}

func generateCommitMessage(userMessage, diffOutput string, stats []fileChangeStats) string {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load config, using simple message generator: %v\n", err)
	}

	llmClient, err := llm.NewLLMClient(cfg)
	if err == nil {
		fmt.Printf("Generating commit message with %s...\n", cfg.Provider)
		llmMessage, err := llmClient.GenerateCommitMessage(diffOutput, userMessage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: LLM message generation failed, falling back to simple generator: %v\n", err)
			return generateSimpleCommitMessage(userMessage, stats)
		}
		return llmMessage
	}
	return generateSimpleCommitMessage(userMessage, stats)
}

func generateSimpleCommitMessage(userMessage string, stats []fileChangeStats) string {
	commitTitle := userMessage
	commitBody := ""

	userMessageLines := strings.SplitN(userMessage, "\n", 2)
	if len(userMessageLines) > 0 {
		commitTitle = userMessageLines[0]
		if len(userMessageLines) > 1 {
			commitBody = userMessageLines[1]
		}
	}

	var b strings.Builder
	b.WriteString(commitTitle + "\n\n")
	if commitBody != "" {
		b.WriteString(commitBody + "\n\n")
	}

	if len(stats) > 0 {
		b.WriteString("Changes:\n")
		for _, stat := range stats {
			b.WriteString(fmt.Sprintf("- %s (%d lines, %d characters)\n", stat.filePath, stat.linesChanged, stat.charsChanged))
		}
	} else {
		b.WriteString("No specific file changes detected in diff.\n")
	}

	return b.String()
}
