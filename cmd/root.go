package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gmh5225/codelens/pkg/version"
	"github.com/spf13/cobra"
)

var (
	// CLI flags
	maxFileSize    float64
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

		// Validate max-size
		if maxFileSize == 0 {
			return fmt.Errorf("--max-size cannot be 0. Use -1 for no limit or a positive number for MB limit")
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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of CodeLens",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("CodeLens %s\n", version.Version)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// check if the error message is related to the repository
func isRepoError(errMsg string) bool {
	repoErrors := []string{
		"repository not found",
		"clone failed",
		"authentication failed",
		"repository does not exist",
		"remote: Repository not found",
		"fatal: repository",
		"could not read from remote repository",
		"exit status",
	}

	errMsg = strings.ToLower(errMsg)
	for _, repoErr := range repoErrors {
		if strings.Contains(errMsg, strings.ToLower(repoErr)) {
			return true
		}
	}
	return false
}

func init() {
	// set custom error handling
	rootCmd.SilenceUsage = true  // default not show usage
	rootCmd.SilenceErrors = true // let us handle the error output

	// add error handling function
	cobra.OnFinalize(func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				if isRepoError(e.Error()) {
					fmt.Fprintf(os.Stderr, "Error: %v\n", e)
				} else {
					fmt.Fprintf(os.Stderr, "Error: %v\n", e)
					rootCmd.Usage()
				}
				os.Exit(1)
			}
			panic(err) // rethrow non-error type panic
		}
	})

	rootCmd.AddCommand(versionCmd)

	// Define flags
	rootCmd.PersistentFlags().Float64VarP(&maxFileSize, "max-size", "s", 10, "Maximum size per file in MB (default: 10MB). Files larger than this will be skipped during analysis")
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
