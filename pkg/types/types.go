package types

// CloneConfig configuration for repository cloning
type CloneConfig struct {
	URL         string // Repository URL
	LocalPath   string // Local storage path
	Branch      string // Branch name (optional)
	Commit      string // Commit hash (optional)
	KeepFiles   bool   // Whether to keep cloned files
	ValidateURL bool
	// New options
	FilterPaths []string // Only clone specific paths
	Depth       int      // Clone depth (0 for full history)
	NoTags      bool     // Don't clone tags
}

// CodeLensConfig configuration for code collection
type CodeLensConfig struct {
	Path            string   // Local code path
	MaxFileSize     int64    // Maximum file size limit in bytes (internally converted)
	IncludePatterns []string // File patterns to include
	ExcludePatterns []string // File patterns to exclude
}

// FileContent represents collected file content
type FileContent struct {
	Path      string // File path
	Content   string // File content
	Size      int64  // File size in bytes
	FileType  string // File type/extension
	LineCount int    // Total number of lines in the file
}

// CodeLensResult represents collection results
type CodeLensResult struct {
	Files      []FileContent            // Collected files
	TotalSize  int64                    // Total size
	TotalFiles int                      // Total number of files
	Summary    string                   // Summary information
	Languages  map[string]LanguageStats // Language statistics
}

// LanguageStats represents statistics for a programming language
type LanguageStats struct {
	Files int64  // Number of files
	Size  int64  // Total size in bytes
	Color string // Language color (for visualization)
}
