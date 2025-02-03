package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var cleanallCmd = &cobra.Command{
	Use:   "cleanall",
	Short: "Clean all cached repositories",
	Long:  `Remove all cached repositories in ~/.codelens directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}

		// Get .codelens directory
		codelensDir := filepath.Join(homeDir, ".codelens")

		// Check if directory exists
		if _, err := os.Stat(codelensDir); os.IsNotExist(err) {
			fmt.Println("Cache directory does not exist. Nothing to clean.")
			return nil
		}

		// Remove the entire .codelens directory
		fmt.Println("Cleaning all cached repositories...")
		if err := os.RemoveAll(codelensDir); err != nil {
			return fmt.Errorf("failed to clean cache directory: %w", err)
		}

		fmt.Println("All cached repositories have been removed successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanallCmd)
}
