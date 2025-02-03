# CodeLens

CodeLens is a command-line tool for analyzing source code repositories. It can analyze both local directories and GitHub repositories, generating comprehensive markdown reports suitable for LLM (Large Language Model) processing.

## Features

- Analyze GitHub repositories directly
- Process local code directories
- Configurable file size limits
- Customizable file inclusion/exclusion patterns
- Markdown report generation
- Git repository caching
- Support for multiple programming languages

## Installation

```bash
go install github.com/gmh5225/codelens@main
```

This will install the `codelens` binary to your `$GOPATH/bin` directory. Make sure your `$GOPATH/bin` is in your `PATH`.

Alternatively, you can specify a specific version or commit:
```bash
go install github.com/gmh5225/codelens@v1.0.1
# or use a specific commit
go install github.com/gmh5225/codelens@commit-hash
```

## Usage

> Note: The `--repo` and `--path` flags are mutually exclusive. You must specify exactly one of them.

### Analyze GitHub Repository
```bash
codelens --repo https://github.com/SimonWaldherr/golang-examples --output ./analysis

# Analyze specific branch
codelens --repo https://github.com/SimonWaldherr/golang-examples --branch master --output ./analysis
```

### Analyze Local Directory
```bash
codelens --path ./myproject --output ./analysis
```

### Advanced Options
```bash
# Set maximum file size (e.g., 5MB)
codelens --path ./myproject --max-size 5

# Specify file patterns
codelens --path ./myproject \
  --include "*.go" --include "*.md" \
  --exclude "vendor/*" --exclude "*_test.go"

# Clone with specific depth and patterns
codelens --repo https://github.com/SimonWaldherr/golang-examples \
  --depth 1 \
  --filter "*.go" --filter "*.md" \
  --skip-tags=false

# Clean up repository after analysis
codelens --repo https://github.com/SimonWaldherr/golang-examples --clean

# Validate Git repository URL(for security)
codelens --repo https://github.com/SimonWaldherr/golang-examples --validate

# Clean all cached repositories
codelens cleanall
```
## Command Line Options

### Global Flags

- `--repo, -r`: GitHub repository URL
- `--path, -p`: Local directory path
- `--output, -o`: Output directory (default: ".")
- `--branch, -b`: Git branch to analyze (default: repository's default branch)
- `--clean, -c`: Clean up repository after analysis
- `--validate, -v`: Validate Git repository URL
- `--depth, -d`: Git clone depth (default: 1, 0 for full history)
- `--skip-tags`: Skip downloading Git tags (default: true)
- `--filter, -f`: File patterns to clone (empty for all files)
- `--max-size, -s`: Maximum file size in MB (-1 for no limit, cannot be 0)
- `--include, -i`: File patterns to include (empty for all files)
- `--exclude, -e`: File patterns to exclude (empty for no exclusions)

### Commands

- `cleanall`: Remove all cached repositories in ~/.codelens directory

## Output Format

CodeLens generates a markdown report (`codelens.md`) with the following structure:

```markdown
# Source Code Analysis for Repository: <repo-name>

This document contains a comprehensive analysis of the source code...

## Analysis Configuration
- Max file size: 10.00 MB
- Include patterns: [*.go *.md]
- Exclude patterns: [vendor/* *_test.go]

## Repository Overview
Key statistics about the analyzed codebase...

## File Structure
Below is the list of analyzed source files...

## File Contents
......
```

## Cache Directory

CodeLens caches GitHub repositories in `~/.codelens/repos/` to avoid repeated downloads.

## Requirements

- Go 1.16 or later
- Git (for repository cloning)

## Development

### Code Formatting

To format all Go files in the repository:
```bash
# Using gofmt (built-in)
gofmt -w .
```

The project follows standard Go formatting conventions. It's recommended to format your code before submitting any changes.

