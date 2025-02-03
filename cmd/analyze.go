package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	// Check if repository already exists
	if _, err := os.Stat(repoDir); err == nil {
		fmt.Println("Using existing repository...")
	} else {
		// Clone repository
		cloneConfig := types.CloneConfig{
			URL:       repoURL,
			LocalPath: repoDir,
			Branch:    branch,
			KeepFiles: !cleanRepo,
		}

		fmt.Println("Cloning repository...")
		if err := git.CloneRepository(cloneConfig); err != nil {
			return fmt.Errorf("clone failed: %w", err)
		}
	}

	// Analyze the repository
	err = analyzeLocalPath(repoDir, outputDir)

	// Clean up if requested and if we cloned the repo
	if cleanRepo {
		fmt.Println("Cleaning up repository...")
		if _, err := os.Stat(repoDir); err == nil {
			if cleanErr := os.RemoveAll(repoDir); cleanErr != nil {
				fmt.Printf("Warning: failed to clean up repository: %v\n", cleanErr)
			} else {
				fmt.Println("Repository cleaned up successfully")
			}
		}
	}

	return err
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
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Extract repository name from URL
	repoName := filepath.Base(strings.TrimSuffix(repoURL, ".git"))

	// Return full path including repository name
	return filepath.Join(cacheDir, repoName), nil
}
