package codelens

import (
	"path/filepath"
	"strings"
)

// getFileType determines the file type based on extension
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
