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

// CloneRepository clones a Git repository
func CloneRepository(config types.CloneConfig) error {
	// Check if git is installed
	if err := checkGitInstalled(); err != nil {
		return err
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
		args = append(args, "--filter=blob:none")
	}

	// Add depth parameter
	if config.Depth > 0 {
		args = append(args, fmt.Sprintf("--depth=%d", config.Depth))
	}

	// Add branch if specified
	if config.Branch != "" {
		args = append(args, "--branch", config.Branch)
	}

	// Add single branch option
	args = append(args, "--single-branch")

	// Skip tags if requested
	if config.NoTags {
		args = append(args, "--no-tags")
	}

	// Add repository URL and local path
	args = append(args, config.URL, config.LocalPath)

	fmt.Printf("Running git command: git %s\n", strings.Join(args, " "))
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	// If using sparse checkout, set up patterns
	if len(config.FilterPaths) > 0 && config.FilterPaths[0] != "" {
		// First enable sparse-checkout
		cmd = exec.Command("git", "-C", config.LocalPath, "config", "core.sparseCheckout", "true")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to enable sparse-checkout: %w", err)
		}

		// Create sparse-checkout file with patterns
		sparseFile := filepath.Join(config.LocalPath, ".git", "info", "sparse-checkout")
		content := strings.Join(config.FilterPaths, "\n")
		if err := os.WriteFile(sparseFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write sparse-checkout patterns: %w", err)
		}

		// Update the working tree
		cmd = exec.Command("git", "-C", config.LocalPath, "read-tree", "-mu", "HEAD")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to update working tree: %w", err)
		}
	}

	// If commit is specified, perform checkout
	if config.Commit != "" {
		cmd = exec.Command("git", "-C", config.LocalPath, "checkout", config.Commit)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
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

// setupSparseCheckout configures sparse checkout for specific paths
func setupSparseCheckout(repoPath string, paths []string) error {
	// Enable sparse-checkout
	cmd := exec.Command("git", "-C", repoPath, "config", "core.sparseCheckout", "true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Create sparse-checkout file
	sparseFile := filepath.Join(repoPath, ".git", "info", "sparse-checkout")
	if err := os.MkdirAll(filepath.Dir(sparseFile), 0755); err != nil {
		return fmt.Errorf("failed to create sparse-checkout directory: %w", err)
	}

	// Write patterns to sparse-checkout file
	content := strings.Join(paths, "\n")
	if err := os.WriteFile(sparseFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write sparse-checkout file: %w", err)
	}

	// Read the working tree
	cmd = exec.Command("git", "-C", repoPath, "read-tree", "-mu", "HEAD")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to read-tree: %w", err)
	}

	return nil
}
