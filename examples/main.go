package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "sort"
    
    "github.com/gmh5225/GodeLens/pkg/git"
    "github.com/gmh5225/GodeLens/pkg/codelens"
    "github.com/gmh5225/GodeLens/pkg/types"
)

// getRepoDir gets the local cache directory for the repository
func getRepoDir(repoURL string) (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("failed to get user home directory: %w", err)
    }
    
    // Create .godelens/repos directory under user's home
    cacheDir := filepath.Join(homeDir, ".godelens", "repos")
    if err := os.MkdirAll(cacheDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create cache directory: %w", err)
    }
    
    // Use repository name as subdirectory
    repoName := filepath.Base(strings.TrimSuffix(repoURL, ".git"))
    return filepath.Join(cacheDir, repoName), nil
}

// generateDirTree generates a tree structure of the directory
func generateDirTree(files []types.FileContent, baseDir string) string {
    var tree strings.Builder
    tree.WriteString("Directory Structure:\n")

    // Get all relative paths and sort them
    var paths []string
    for _, file := range files {
        relPath, err := filepath.Rel(baseDir, file.Path)
        if err == nil {
            paths = append(paths, relPath)
        }
    }
    sort.Strings(paths)

    // Generate tree structure
    for i, path := range paths {
        parts := strings.Split(path, string(os.PathSeparator))
        indent := ""
        for j := 0; j < len(parts)-1; j++ {
            if j == len(parts)-2 {
                indent += "├── "
            } else {
                indent += "│   "
            }
        }
        // Use different symbol for the last file
        if i == len(paths)-1 {
            indent = strings.ReplaceAll(indent, "├──", "└──")
        }
        tree.WriteString(indent + parts[len(parts)-1] + "\n")
    }
    return tree.String()
}

func main() {
    repoURL := "https://github.com/SimonWaldherr/golang-examples"
    repoDir, err := getRepoDir(repoURL)
    if err != nil {
        log.Fatal(err)
    }

    // Clone repository if directory doesn't exist
    if _, err := os.Stat(repoDir); os.IsNotExist(err) {
        // Clone example repository
        cloneConfig := types.CloneConfig{
            URL:       repoURL,
            LocalPath: repoDir,
            Branch:    "master",
            KeepFiles: true,  // Keep cloned files
        }
        
        fmt.Println("Cloning repository...")
        if err := git.CloneRepository(cloneConfig); err != nil {
            log.Fatal("Clone failed:", err)
        }
        fmt.Println("Clone completed")
    } else {
        fmt.Println("Using cached repository")
    }

    // Collect code
    config := types.CodeLensConfig{
        Path:        repoDir,
        MaxFileSize: 10 * 1024 * 1024, // 10MB
        IncludePatterns: []string{
            "*.go",    // Only Go files
            "*.md",    // And Markdown files
        },
        ExcludePatterns: []string{
            "vendor/*",   // Exclude vendor directory
            "*_test.go",  // Exclude test files
        },
    }
    
    fmt.Println("Starting code collection...")
    result, err := codelens.CollectCode(config)
    if err != nil {
        log.Fatal("Code collection failed:", err)
    }

    // Print results
    fmt.Println("\n=== Collection Results ===")
    fmt.Println(result.Summary)
    fmt.Println("\n=== Directory Structure ===")
    fmt.Println(generateDirTree(result.Files, repoDir))
    fmt.Println("\n=== File Contents ===")
    
    // Prepare content for file writing
    var content strings.Builder
    
    // Write basic information
    content.WriteString(fmt.Sprintf("# Repository Code Analysis: %s\n\n", filepath.Base(repoDir)))
    content.WriteString("This is an automated code analysis of the repository, including file structure and source code content.\n\n")
    content.WriteString("## Summary\n")
    content.WriteString(result.Summary)
    content.WriteString("\n")
    
    // Write simplified directory structure
    var paths []string
    for _, file := range result.Files {
        relPath, _ := filepath.Rel(repoDir, file.Path)
        paths = append(paths, relPath)
    }
    sort.Strings(paths)
    content.WriteString("## Repository Structure\n")
    content.WriteString("The following files were analyzed:\n\n")
    for _, path := range paths {
        content.WriteString("- " + path + "\n")
    }
    
    content.WriteString("\n")
    
    for _, file := range result.Files {
        // Get relative path
        relPath, err := filepath.Rel(repoDir, file.Path)
        if err != nil {
            log.Printf("Warning: Unable to get relative path %s: %v", file.Path, err)
            relPath = file.Path
        }
        
        // Add to file content
        content.WriteString(fmt.Sprintf("## %s\n", relPath))
        content.WriteString("```" + file.FileType + "\n")
        content.WriteString(file.Content)
        content.WriteString("\n```\n\n")
    }
    
    // Save to file
    outputPath := filepath.Join(repoDir, "codelens.md")
    // Remove file if it exists
    if err := os.Remove(outputPath); err != nil && !os.IsNotExist(err) {
        log.Fatal("Failed to delete old file:", err)
    }
    if err := os.WriteFile(outputPath, []byte(content.String()), 0644); err != nil {
        log.Fatal("Failed to save file:", err)
    }
    
    fmt.Printf("\nResults saved to: %s\n", outputPath)
} 