package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gmh5225/codelens/pkg/types"
)

// checkGitInstalled verifies if git is available in the system
func checkGitInstalled() error {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git is not installed or not available in PATH: %w", err)
	}
	return nil
}

// runGitCommand executes a git command with the given args
func runGitCommand(repoPath string, args ...string) error {
	cmd := exec.Command("git", args...)
	if repoPath != "" {
		args = append([]string{"-C", repoPath}, args...)
		cmd = exec.Command("git", args...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// setupSparseCheckout configures sparse-checkout with given patterns
func setupSparseCheckout(repoPath string, patterns []string) error {
	// Enable sparse-checkout without cone mode
	if err := runGitCommand(repoPath, "sparse-checkout", "init", "--no-cone"); err != nil {
		return fmt.Errorf("failed to enable sparse-checkout: %w", err)
	}

	// Set patterns
	if err := runGitCommand(repoPath, "sparse-checkout", "set", "--no-cone", strings.Join(patterns, " ")); err != nil {
		return fmt.Errorf("failed to set sparse-checkout patterns: %w", err)
	}

	return nil
}

// CloneRepository clones a Git repository
func CloneRepository(config types.CloneConfig) error {
	// Check if git is installed
	if err := runGitCommand("", "--version"); err != nil {
		return fmt.Errorf("git is not installed or not available in PATH: %w", err)
	}

	// Check if repository already exists
	if _, err := os.Stat(config.LocalPath); err == nil {
		fmt.Println("Repository exists, updating...")

		// Fetch latest changes
		if err := runGitCommand(config.LocalPath, "fetch", "origin"); err != nil {
			return fmt.Errorf("failed to fetch updates: %w", err)
		}

		// Reset to latest changes
		if err := runGitCommand(config.LocalPath, "reset", "--hard", "origin/HEAD"); err != nil {
			return fmt.Errorf("failed to reset to latest changes: %w", err)
		}

		// Configure sparse-checkout if needed
		if len(config.FilterPaths) > 0 && config.FilterPaths[0] != "" {
			if err := setupSparseCheckout(config.LocalPath, config.FilterPaths); err != nil {
				return err
			}
		} else {
			// Disable sparse-checkout if no filters
			if err := runGitCommand(config.LocalPath, "sparse-checkout", "disable"); err != nil {
				return fmt.Errorf("failed to disable sparse-checkout: %w", err)
			}

			// Checkout all files
			if err := runGitCommand(config.LocalPath, "checkout", "HEAD"); err != nil {
				return fmt.Errorf("failed to checkout files: %w", err)
			}
		}

		return nil
	}

	// Validate repository URL if required
	if config.ValidateURL {
		if err := validateGitURL(config.URL); err != nil {
			return fmt.Errorf("invalid repository URL: %w", err)
		}
	}

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(config.LocalPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	args := []string{"clone"}

	// Add sparse checkout if paths are specified
	if len(config.FilterPaths) > 0 && config.FilterPaths[0] != "" {
		args = append(args, "--sparse")
		args = append(args, "--no-checkout") // 先不检出文件，等设置好 sparse-checkout 后再检出
	}

	// Add depth parameter
	if config.Depth > 0 {
		args = append(args, fmt.Sprintf("--depth=%d", config.Depth))
	}

	// Add branch if specified
	if config.Branch != "" {
		args = append(args, "--branch", config.Branch)
	}

	// Skip tags if requested
	if config.NoTags {
		args = append(args, "--no-tags")
	}

	// Add repository URL and local path
	args = append(args, config.URL, config.LocalPath)

	fmt.Printf("Running git command: git %s\n", strings.Join(args, " "))
	if err := runGitCommand("", args...); err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	// If using sparse checkout, set up patterns
	if len(config.FilterPaths) > 0 && config.FilterPaths[0] != "" {
		if err := setupSparseCheckout(config.LocalPath, config.FilterPaths); err != nil {
			return err
		}

		// Checkout files
		if err := runGitCommand(config.LocalPath, "checkout"); err != nil {
			return fmt.Errorf("failed to checkout files: %w", err)
		}
	}

	// If commit is specified, perform checkout
	if config.Commit != "" {
		if err := runGitCommand(config.LocalPath, "checkout", config.Commit); err != nil {
			return fmt.Errorf("checkout failed: %w", err)
		}
	}

	return nil
}

func validateGitURL(url string) error {
	if !IsValidGitHost(url) {
		return fmt.Errorf("unsupported Git host. Supported hosts: %v", SupportedGitHosts)
	}

	// Validate URL format
	parts := strings.Split(strings.TrimSuffix(url, ".git"), "/")
	if len(parts) < 5 {
		return fmt.Errorf("invalid repository URL format")
	}

	// Validate username and repository name
	if parts[len(parts)-2] == "" || parts[len(parts)-1] == "" {
		return fmt.Errorf("invalid username or repository name")
	}

	return nil
}
