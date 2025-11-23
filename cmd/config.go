package cmd

import (
	"fmt"
	"github.com/biswajitpain/gitter/internal/config"


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
gitter config --provider openai --api-key sk-...
gitter config --provider gemini --api-key your-gemini-api-key`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if provider == "" && apiKey == "" {
			cmd.Help()
			return nil // Showing help is not an error, return nil
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("error loading existing config: %w", err)
		}

		if provider != "" {
			cfg.Provider = provider
		}
		if apiKey != "" {
			cfg.APIKey = apiKey
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		fmt.Println("Configuration updated successfully.")
		if cfg.Provider != "" {
			fmt.Printf("Provider: %s\n", cfg.Provider)
		}
		if cfg.APIKey != "" {
			fmt.Println("API Key: [set]")
		}
		return nil // Success, return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&provider, "provider", "p", "", "The LLM provider (e.g., 'openai', 'gemini')")
	configCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "The API key for the LLM provider")
}
