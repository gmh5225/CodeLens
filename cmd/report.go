package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gmh5225/codelens/pkg/git"
	"github.com/gmh5225/codelens/pkg/types"
)

func generateReport(result *types.CodeLensResult, basePath, outputPath string) error {
	var content strings.Builder

	// Check if there are any files
	if len(result.Files) == 0 {
		content.WriteString("# No files were analyzed\n\n")
		content.WriteString("Possible reasons:\n")
		content.WriteString("- All files exceeded size limit\n")
		content.WriteString("- No files matched include patterns\n")
		content.WriteString("- All files matched exclude patterns\n")
		return os.WriteFile(outputPath, []byte(content.String()), 0644)
	}

	// Write header
	content.WriteString(fmt.Sprintf("# Source Code Analysis for Repository: %s\n\n", filepath.Base(basePath)))
	content.WriteString("This document contains a comprehensive analysis of the source code, including file structure and content. The analysis is designed to help understand the codebase structure and implementation details.\n\n")

	content.WriteString("## Analysis Configuration\n")
	if maxFileSize == -1 {
		content.WriteString("- Max file size: no limit\n")
	} else {
		content.WriteString(fmt.Sprintf("- Max file size: %.2f MB\n", float64(maxFileSize)))
	}
	content.WriteString(fmt.Sprintf("- Include patterns: %v\n", includePattern))
	content.WriteString(fmt.Sprintf("- Exclude patterns: %v\n", excludePattern))
	content.WriteString("\n")

	// Write language statistics if available
	if githubRepo != "" {
		langInfo, err := git.GetRepoLanguages(githubRepo)
		if err == nil && langInfo.PrimaryLang != "" {
			content.WriteString("## Language Statistics\n")
			content.WriteString("Based on GitHub's language detection:\n\n")
			content.WriteString(fmt.Sprintf("Primary language: **%s** (%.2f%%)\n\n",
				langInfo.PrimaryLang,
				langInfo.GetPercentage(langInfo.PrimaryLang),
			))
		}
	}

	// Write summary
	content.WriteString("## Repository Overview\n")
	content.WriteString("Key statistics about the analyzed codebase:\n\n")
	content.WriteString(result.Summary)

	// Write file structure
	content.WriteString("## File Structure\n")
	content.WriteString("Below is the list of analyzed source files in this repository:\n\n")

	var paths []string
	for _, file := range result.Files {
		relPath, _ := filepath.Rel(basePath, file.Path)
		paths = append(paths, relPath)
	}
	sort.Strings(paths)

	for _, path := range paths {
		content.WriteString("- " + path + "\n")
	}
	content.WriteString("\n")

	// Write file contents
	content.WriteString("## File Contents\n\n")
	for _, file := range result.Files {
		relPath, _ := filepath.Rel(basePath, file.Path)
		content.WriteString(fmt.Sprintf("### %s (%d lines)\n", relPath, file.LineCount))
		content.WriteString("```" + file.FileType + "\n")
		content.WriteString(file.Content)
		content.WriteString("\n```\n\n")
	}

	return os.WriteFile(outputPath, []byte(content.String()), 0644)
}
