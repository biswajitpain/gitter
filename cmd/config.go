package cmd

import (
	"fmt"
	"gitter/internal/config"
	"os"

	"github.com/spf13/cobra"
)

var (
	provider string
	apiKey   string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the LLM provider and API key",
	Long: `Configure the settings for the LLM provider used for generating commit messages.
Currently supports OpenAI.

Example:
gitter config --provider openai --api-key sk-...`,
	Run: func(cmd *cobra.Command, args []string) {
		if provider == "" && apiKey == "" {
			cmd.Help()
			os.Exit(0)
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading existing config: %v\n", err)
			os.Exit(1)
		}

		if provider != "" {
			cfg.Provider = provider
		}
		if apiKey != "" {
			cfg.APIKey = apiKey
		}

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Configuration updated successfully.")
		if cfg.Provider != "" {
			fmt.Printf("Provider: %s\n", cfg.Provider)
		}
		if cfg.APIKey != "" {
			fmt.Println("API Key: [set]")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&provider, "provider", "p", "", "The LLM provider (e.g., 'openai')")
	configCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "The API key for the LLM provider")
}
