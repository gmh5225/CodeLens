package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gmh5225/codelens/pkg/codelens"
	"github.com/gmh5225/codelens/pkg/git"
	"github.com/gmh5225/codelens/pkg/types"
)

func analyzeGitRepo(repoURL, outputDir string) error {
	// Get repo directory
	repoDir, err := getRepoDir(repoURL)
	if err != nil {
		return err
	}

	// Clone repository
	cloneConfig := types.CloneConfig{
		URL:       repoURL,
		LocalPath: repoDir,
		Branch:    branch,
		KeepFiles: true,
	}

	fmt.Println("Cloning repository...")
	if err := git.CloneRepository(cloneConfig); err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	return analyzeLocalPath(repoDir, outputDir)
}

func analyzeLocalPath(path, outputDir string) error {
	fmt.Println("Starting code analysis...")

	config := types.CodeLensConfig{
		Path:            path,
		MaxFileSize:     maxFileSize,
		IncludePatterns: includePattern,
		ExcludePatterns: excludePattern,
	}

	result, err := codelens.CollectCode(config)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Generate output file path
	outputFile := filepath.Join(outputDir, "codelens.md")
	if err := generateReport(result, path, outputFile); err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	fmt.Printf("\nAnalysis complete. Results saved to: %s\n", outputFile)
	return nil
}

func getRepoDir(repoURL string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create .codelens/repos directory under user's home
	cacheDir := filepath.Join(homeDir, ".codelens", "repos")

	return cacheDir, nil
}
