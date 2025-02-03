package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(config.LocalPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	args := []string{"clone"}

	// If branch is specified
	if config.Branch != "" {
		args = append(args, "--branch", config.Branch)
	}

	// Add shallow clone parameters
	args = append(args, "--depth=1", "--single-branch")

	// Add repository URL and local path
	args = append(args, config.URL, config.LocalPath)

	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("clone failed: %s: %w", output, err)
	}

	// If commit is specified, perform checkout
	if config.Commit != "" {
		checkoutCmd := exec.Command("git", "-C", config.LocalPath, "checkout", config.Commit)
		if output, err := checkoutCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("checkout failed: %s: %w", output, err)
		}
	}

	return nil
}
