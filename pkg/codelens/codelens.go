package codelens

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gmh5225/codelens/pkg/types"
)

// DefaultIgnorePatterns default patterns for files to ignore
var DefaultIgnorePatterns = []string{
	// Version control
	".git", ".svn", ".hg",
	".gitignore", ".gitattributes", ".gitmodules",

	// Build and dependency directories
	"node_modules", "vendor", "dist", "build",
	"target", "bin", "obj",

	// IDEs and editors
	".idea", ".vscode", ".vs",
	"*.swp", "*.swo", "*.swn",

	// Temporary files
	"*.tmp", "*.temp", "*.bak",
	".DS_Store", "Thumbs.db",

	// Binary and executable files
	"*.exe", "*.dll", "*.so", "*.dylib",
	"*.bin", "*.dat",

	// Compressed files
	"*.zip", "*.rar", "*.7z", "*.gz", "*.tar",

	// Media files
	"*.jpg", "*.jpeg", "*.png", "*.gif", "*.ico",
	"*.mp3", "*.mp4", "*.avi", "*.mov",

	// Compilation cache
	"*.pyc", "*.pyo", "*.pyd", "__pycache__",
	"*.class", "*.o", "*.obj",
}

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
		content, err := readFileContent(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Add to results
		result.Files = append(result.Files, types.FileContent{
			Path:     path,
			Content:  content,
			Size:     info.Size(),
			FileType: getFileType(path),
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

func readFileContent(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getFileType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	// Programming languages
	case ".go":
		return "golang"
	case ".py":
		return "python"
	case ".js", ".jsx", ".ts", ".tsx":
		return "javascript"
	case ".java":
		return "java"
	case ".cpp", ".cc", ".cxx", ".c++", ".hpp":
		return "cpp"
	case ".c", ".h":
		return "c"
	case ".cs":
		return "csharp"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".rs":
		return "rust"
	case ".swift":
		return "swift"
	case ".kt", ".kts":
		return "kotlin"

	// Markup languages and configuration files
	case ".md", ".markdown":
		return "markdown"
	case ".json":
		return "json"
	case ".xml":
		return "xml"
	case ".yaml", ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".ini":
		return "ini"

	// Scripting languages
	case ".sh", ".bash":
		return "shell"
	case ".ps1":
		return "powershell"
	case ".bat", ".cmd":
		return "batch"

	// Other common file types
	case ".sql":
		return "sql"
	case ".html", ".htm":
		return "html"
	case ".css", ".scss", ".sass", ".less":
		return "css"

	case "":
		return "unknown"
	default:
		// Return extension without dot
		return ext[1:]
	}
}
