package types

// CloneConfig configuration for repository cloning
type CloneConfig struct {
	URL       string // Repository URL
	LocalPath string // Local storage path
	Branch    string // Branch name (optional)
	Commit    string // Commit hash (optional)
	KeepFiles bool   // Whether to keep cloned files
}

// CodeLensConfig configuration for code collection
type CodeLensConfig struct {
	Path            string   // Local code path
	MaxFileSize     int64    // Maximum file size limit
	IncludePatterns []string // File patterns to include
	ExcludePatterns []string // File patterns to exclude
}

// FileContent represents collected file content
type FileContent struct {
	Path     string // File path
	Content  string // File content
	Size     int64  // File size
	FileType string // File type
}

// CodeLensResult represents collection results
type CodeLensResult struct {
	Files      []FileContent // Collected files
	TotalSize  int64         // Total size
	TotalFiles int           // Total number of files
	Summary    string        // Summary information
}
