package codelens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gmh5225/codelens/pkg/types"
)

// CollectCode collects code from local path
func CollectCode(config types.CodeLensConfig) (*types.CodeLensResult, error) {
	result := &types.CodeLensResult{
		Files: make([]types.FileContent, 0),
	}

	// Merge default and custom ignore patterns
	excludePatterns := append(DefaultIgnorePatterns, config.ExcludePatterns...)

	err := filepath.Walk(config.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			if shouldSkipDir(path, excludePatterns) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check file size
		if config.MaxFileSize != -1 && info.Size() > config.MaxFileSize {
			return nil
		}

		// Check if file should be included
		if !shouldIncludeFile(path, config.IncludePatterns, excludePatterns) {
			return nil
		}

		// Read file content
		content, lineCount, err := readFileContent(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Add to results
		result.Files = append(result.Files, types.FileContent{
			Path:      path,
			Content:   content,
			Size:      info.Size(),
			FileType:  getFileType(path),
			LineCount: lineCount,
		})

		result.TotalSize += info.Size()
		result.TotalFiles++

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to traverse directory: %w", err)
	}

	// Generate summary information
	result.Summary = fmt.Sprintf(
		"# Summary\n- Total files: %d\n- Total size: %.2f MB\n",
		result.TotalFiles,
		float64(result.TotalSize)/(1024*1024),
	)

	return result, nil
}

// Helper functions
func shouldSkipDir(path string, excludePatterns []string) bool {
	base := filepath.Base(path)
	for _, pattern := range excludePatterns {
		if matched, _ := filepath.Match(pattern, base); matched {
			return true
		}
	}
	return false
}

func shouldIncludeFile(path string, includePatterns, excludePatterns []string) bool {
	// If no include patterns are specified, include all files
	if len(includePatterns) == 0 {
		return !isExcluded(path, excludePatterns)
	}

	// Check if file matches include patterns
	for _, pattern := range includePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return !isExcluded(path, excludePatterns)
		}
	}

	return false
}

func isExcluded(path string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}

func readFileContent(path string) (string, int, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", 0, err
	}

	// Convert content to string and count lines
	contentStr := string(content)
	lineCount := strings.Count(contentStr, "\n") + 1

	return contentStr, lineCount, nil
}
