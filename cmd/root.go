package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// CLI flags
	maxFileSize    int64
	filterPaths    []string
	localPath      string
	githubRepo     string
	outputDir      string
	branch         string
	cleanRepo      bool
	validateURL    bool
	cloneDepth     int
	skipTags       bool
	includePattern []string
	excludePattern []string
)

var rootCmd = &cobra.Command{
	Use:   "codelens",
	Short: "CodeLens is a code analysis tool",
	Long: `CodeLens analyzes code repositories and generates a markdown report.
It can analyze both local directories and GitHub repositories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate flags
		if githubRepo != "" && localPath != "" {
			return fmt.Errorf("cannot specify both --repo and --path flags")
		}

		if githubRepo == "" && localPath == "" {
			return fmt.Errorf("must specify either --repo or --path flag")
		}

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		if githubRepo != "" {
			return analyzeGitRepo(githubRepo, outputDir)
		}
		return analyzeLocalPath(localPath, outputDir)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Define flags
	rootCmd.PersistentFlags().Int64VarP(&maxFileSize, "max-size", "s", 10*1024*1024, "Maximum file size in bytes")
	rootCmd.PersistentFlags().StringVarP(&localPath, "path", "p", "", "Local path to analyze")
	rootCmd.PersistentFlags().StringVarP(&githubRepo, "repo", "r", "", "GitHub repository URL")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", ".", "Output directory for analysis results")
	rootCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "Git branch to analyze (default: repository's default branch)")
	rootCmd.PersistentFlags().BoolVarP(&cleanRepo, "clean", "c", false, "Clean up repository after analysis")
	rootCmd.PersistentFlags().BoolVarP(&validateURL, "validate", "v", false, "Validate Git repository URL")
	rootCmd.PersistentFlags().IntVarP(&cloneDepth, "depth", "d", 1, "Git clone depth (0 for full history)")
	rootCmd.PersistentFlags().BoolVar(&skipTags, "skip-tags", true, "Skip downloading Git tags")
	rootCmd.PersistentFlags().StringArrayVarP(&filterPaths, "filter", "f", []string{}, "File patterns to clone (empty for all files)")
	rootCmd.PersistentFlags().StringArrayVarP(&includePattern, "include", "i", []string{}, "File patterns to include (empty for all files)")
	rootCmd.PersistentFlags().StringArrayVarP(&excludePattern, "exclude", "e", []string{}, "File patterns to exclude (empty for no exclusions)")
}
