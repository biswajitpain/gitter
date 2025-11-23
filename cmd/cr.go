package cmd

import (
	"bufio"
	"fmt"
	"github.com/biswajitpain/gitter/internal/config"
	"github.com/biswajitpain/gitter/internal/llm"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// crCmd represents the cr command

var crCmd = &cobra.Command{

	Use:   "cr",

	Short: "Create a commit with an AI-generated message",

	Long: `The 'cr' command automates the commit process.

It stages files, generates a diff, and uses an LLM (if configured) 

to create a conventional commit message.`,

	RunE: handleCrCommand,

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



func handleCrCommand(cmd *cobra.Command, args []string) error {

	// 1. Check if we are in a git repository.

	gitCheckCmd := execCommand("git", "rev-parse", "--is-inside-work-tree")

	if err := gitCheckCmd.Run(); err != nil {

		return fmt.Errorf("not a git repository: %w", err)

	}



	reader := bufio.NewReader(os.Stdin)



	// 2. Check for staged changes.

	stagedCheckCmd := execCommand("git", "diff", "--cached", "--quiet")

	if stagedCheckCmd.Run() != nil { // Exits with 1 if no staged changes

		fmt.Print("No files are currently staged. Would you like to stage all changed files? (y/n): ")

		stageAllInput, _ := reader.ReadString('\n')

		stageAllInput = strings.TrimSpace(strings.ToLower(stageAllInput))



		if stageAllInput == "y" {

			fmt.Println("Staging all changed files...")

			if err := execCommand("git", "add", ".").Run(); err != nil {

				return fmt.Errorf("error staging changes: %w", err)

			}

		} else {

			fmt.Println("Operation cancelled. No files were staged.")

			return nil

		}

	} else {

		fmt.Println("Working on currently staged changes.")

	}



	// 3. Get the diff of staged changes.

	fmt.Println("Generating diff...")

	diffCmd := execCommand("git", "diff", "--staged")

	diffOutputBytes, err := diffCmd.Output()

	if err != nil {

		return fmt.Errorf("error getting diff: %w", err)

	}

	diffOutput := string(diffOutputBytes)



	if strings.TrimSpace(diffOutput) == "" {

		fmt.Println("No changes to commit.")

		execCommand("git", "reset").Run()

		return nil

	}



		// 4. Parse the diff to get stats.



		stats := parseDiffStats(diffOutput)



	



			// 5. Generate an initial commit message from LLM (if configured).



	



			// This will be used as the base for the suggestion.



	



			var llmGeneratedMessage string



	



			var llmErr error



	



			llmGeneratedMessage, llmErr = generateCommitMessage("", diffOutput) // Pass empty user message for initial generation



	



		



	



			// 6. Generate the file change statistics message. This is always included.



	



			statsMessage := generateSimpleCommitMessage("", stats) // User message is empty for stats generation



	



		



	



				// Combine LLM message (if successful) and stats message.



	



		



	



				var suggestedMessage string



	



		



	



				var messageBody strings.Builder



	



		



	



			



	



		



	



				if llmErr == nil {



	



		



	



					// Format LLM message as an unordered list if it's multi-line



	



		



	



					llmMessageLines := strings.Split(llmGeneratedMessage, "\n")



	



		



	



					for _, line := range llmMessageLines {



	



		



	



						line = strings.TrimSpace(line)



	



		



	



						if line != "" {



	



		



	



							messageBody.WriteString("- " + line + "\n")



	



		



	



						}



	



		



	



					}



	



		



	



				} else {



	



		



	



					// If LLM generation failed, use a simple default message and log the error



	



		



	



					fmt.Fprintf(os.Stderr, "Warning: LLM commit message generation failed: %v. Using simple default message.\n", llmErr)



	



		



	



					defaultMsg := createDefaultCommitMessage()



	



		



	



					if defaultMsg != "" {



	



		



	



						messageBody.WriteString("- " + defaultMsg + "\n")



	



		



	



					}



	



		



	



				}



	



		



	



			



	



		



	



				// Append stats message, ensuring a blank line separates them if both exist



	



		



	



				if messageBody.Len() > 0 && statsMessage != "" {



	



		



	



					messageBody.WriteString("\n")



	



		



	



				}



	



		



	



				messageBody.WriteString(statsMessage)



	



		



	



			



	



		



	



				suggestedMessage = messageBody.String()



	



		



	



			



	



		



	



				// 7. Ask the user for a commit message, providing the suggestion.



	



		



	



				fmt.Printf("Please enter a commit message (or press Enter to use the suggestion):\n\n")



	



			fmt.Printf("--- Suggested Commit Message ---\n%s\n--------------------------------\n", suggestedMessage)



	



			fmt.Print("> ")



	



			userProvidedMessage, _ := reader.ReadString('\n')



	



			userProvidedMessage = strings.TrimSpace(userProvidedMessage)



	



		



	



			var finalCommitMessage string



	



			if userProvidedMessage == "" {



	



				finalCommitMessage = suggestedMessage



	



				fmt.Println("Using suggested message.")



	



			} else {



	



				finalCommitMessage = userProvidedMessage



	



				fmt.Println("Using user-provided message.")



	



			}



	



		



	



			// 8. Ask for confirmation.



	



			fmt.Println("\n--- Final Commit Message ---")



	



			fmt.Print(finalCommitMessage)



	



			fmt.Println("\n----------------------------")



	



			fmt.Print("Confirm commit with this message? (y/n): ")



	



			confirmInput, _ := reader.ReadString('\n')



	



			confirmInput = strings.TrimSpace(strings.ToLower(confirmInput))



	



		if confirmInput == "y" {



			// 8. Commit.



			fmt.Println("Committing...")



			commitCmd := execCommand("git", "commit", "-m", finalCommitMessage)

		if err := commitCmd.Run(); err != nil {

			return fmt.Errorf("error committing: %w", err)

		}

		fmt.Println("Commit successful.")

	} else {

		fmt.Println("Commit cancelled. Changes are still staged.")

		fmt.Print("Would you like to unstage the changes? (y/n): ")

		unstageInput, _ := reader.ReadString('\n')

		unstageInput = strings.TrimSpace(strings.ToLower(unstageInput))

		if unstageInput == "y" {

			if err := execCommand("git", "reset").Run(); err != nil {

				fmt.Fprintf(os.Stderr, "Error unstaging changes: %v\n", err) // Print, but don't exit if unstage fails

			} else {

				fmt.Println("Changes have been unstaged.")

			}

		}

	}

	return nil

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

		// Find the file path from the "+++ b/" line
		for _, line := range lines {
			if strings.HasPrefix(line, "+++ b/") {
				currentFile = strings.TrimPrefix(line, "+++ b/")
				// Handle cases where the file is deleted (e.g., +++ b/dev/null)
				if currentFile == "/dev/null" {
					currentFile = "" // Or handle as appropriate for deletions
				}
				break
			}
		}
		if currentFile == "" || currentFile == "/dev/null" {
			continue // Skip if no valid file path found or if it's /dev/null
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
	repoPath, err := execCommand("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not determine repo name for default commit message: %v\n", err)
		return "new git commit" // Fallback default
	}
	repoName := filepath.Base(strings.TrimSpace(string(repoPath)))

	dateStr := timeNow().Format("2006-01-02")
	return fmt.Sprintf("chore: new git commit on %s %s", repoName, dateStr)
}

// newLLMClientFunc is a package-level variable to allow mocking llm.NewLLMClient in tests.
var newLLMClientFunc = llm.NewLLMClient

func generateCommitMessage(userMessage, diffOutput string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		// Log the warning to stderr, but don't fail the LLM generation process itself
		// as we now expect a separate fallback/combining logic in handleCrCommand.
		fmt.Fprintf(os.Stderr, "Warning: could not load config for LLM generation: %v\n", err)
	}

	llmClient, err := newLLMClientFunc(cfg)
	if err != nil {
		return "", fmt.Errorf("LLM client could not be initialized (Provider: '%s', Error: %w)", cfg.Provider, err)
	}

	fmt.Printf("Generating commit message with %s...\n", cfg.Provider)
	llmMessage, err := llmClient.GenerateCommitMessage(diffOutput, userMessage)
	if err != nil {
		return "", fmt.Errorf("LLM message generation failed: %w", err)
	}
	return llmMessage, nil
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