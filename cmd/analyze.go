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

	// Get repository language statistics
	langInfo, err := git.GetRepoLanguages(repoURL)
	if err != nil {
		if strings.Contains(err.Error(), "only supported for GitHub") {
			fmt.Println("Note: Language statistics are only available for GitHub repositories")
		} else {
			fmt.Printf("Warning: Failed to get repository language stats: %v\n", err)
		}
	} else if langInfo.PrimaryLang != "" {
		fmt.Printf("Primary language (GitHub): %s (%.2f%%)\n",
			langInfo.PrimaryLang,
			langInfo.GetPercentage(langInfo.PrimaryLang),
		)
	}

	// Clone or update repository
	cloneConfig := types.CloneConfig{
		URL:         repoURL,
		LocalPath:   repoDir,
		Branch:      branch,
		KeepFiles:   !cleanRepo,
		ValidateURL: validateURL,
		Depth:       cloneDepth,
		NoTags:      skipTags,
		FilterPaths: filterPaths,
	}

	fmt.Println("Cloning repository...")
	if err := git.CloneRepository(cloneConfig); err != nil {
		// extract the original error message from git command
		if strings.Contains(err.Error(), "fatal: ") {
			errMsg := err.Error()
			if idx := strings.Index(errMsg, "fatal: "); idx != -1 {
				return fmt.Errorf("%s", strings.TrimSpace(errMsg[idx:]))
			}
		}
		return fmt.Errorf("repository error: %s", err)
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

	// Generate prefix based on path type
	var prefix string
	if localPath != "" {
		// For local paths, use local-dirname format
		prefix = fmt.Sprintf("local-%s", filepath.Base(localPath))
	} else if githubRepo != "" {
		// For GitHub repos, use author/repo format
		parts := strings.Split(strings.TrimSuffix(githubRepo, ".git"), "/")
		if len(parts) >= 2 {
			author := parts[len(parts)-2]
			repoName := parts[len(parts)-1]
			prefix = fmt.Sprintf("%s-%s", author, repoName)
		}
	} else {
		// This shouldn't happen due to validation in rootCmd
		prefix = fmt.Sprintf("unknown-%s", filepath.Base(path))
	}

	config := types.CodeLensConfig{
		Path:            path,
		MaxFileSize:     int64(maxFileSize * 1024 * 1024), // Convert MB to bytes
		IncludePatterns: includePattern,
		ExcludePatterns: excludePattern,
	}

	result, err := codelens.CollectCode(config)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Generate output file path
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s-codelens.md", prefix))

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
