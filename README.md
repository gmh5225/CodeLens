# CodeLens

CodeLens is a command-line tool for analyzing source code repositories. It can analyze both local directories and Remote repositories, generating comprehensive markdown reports suitable for LLM (Large Language Model) processing.

## Features

- Analyze Remote repositories directly
- Process local code directories
- Configurable file size limits
- Customizable file inclusion/exclusion patterns
- Markdown report generation
- Git repository caching
- Support for multiple programming languages
- Token counting optimized for code analysis:
  - Accurate token counting for both code and text
  - Special handling for programming constructs
  - Support for CJK and Unicode characters
  - URL and email address smart tokenization
## Installation

```bash
go clean -modcache
go install github.com/gmh5225/codelens@main
```

This will install the `codelens` binary to your `$GOPATH/bin` directory. Make sure your `$GOPATH/bin` is in your `PATH`.
```bash
# for bash
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# for zsh
echo 'export GOPATH=$HOME/go' >> ~/.zshrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zshrc
source ~/.zshrc
```

Alternatively, you can specify a specific version or commit:
```bash
go install github.com/gmh5225/codelens@v1.0.6
# or use a specific commit
go install github.com/gmh5225/codelens@commit-hash
```

## Usage

> Note: The `--repo` and `--path` flags are mutually exclusive. You must specify exactly one of them.

### Analyze Remote Repository
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
# Set maximum individual file size (default is 10MB per file)
codelens --path ./myproject --max-size 20  # Analyze files up to 20MB each

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

- `--repo, -r`: Remote repository URL
- `--path, -p`: Local directory path
- `--output, -o`: Output directory (default: ".")
- `--branch, -b`: Git branch to analyze (default: repository's default branch)
- `--clean, -c`: Clean up repository after analysis
- `--validate, -v`: Validate Git repository URL
- `--depth, -d`: Git clone depth (default: 1, 0 for full history)
- `--skip-tags`: Skip downloading Git tags (default: true)
- `--filter, -f`: File patterns to clone (empty for all files)
- `--max-size, -s`: Maximum size per file in MB (default: 10MB). Files larger than this will be skipped
- `--include, -i`: File patterns to include (empty for all files)
- `--exclude, -e`: File patterns to exclude (empty for no exclusions)

### Commands

- `cleanall`: Remove all cached repositories in ~/.codelens directory

## Output Format

CodeLens generates a markdown report (`codelens.md`) with the following structure:

### When Files Are Found
```markdown
# Source Code Analysis for Repository: <repo-name>

This document contains a comprehensive analysis of the source code...

## Summary
- Total files: 115
- Total size: 0.24 MB
- Total tokens: 25880

## Analysis Configuration
- Max file size: 10.00 MB
- Include patterns: [*.go *.md]
- Exclude patterns: [vendor/* *_test.go]

## Language Statistics (GitHub repositories only)
Based on GitHub's language detection:

Primary language: **Go** (85.2%)

## Token Statistics
Token distribution by file...

## Repository Overview
Key statistics about the analyzed codebase...

## File Structure
Below is the list of analyzed source files...

## File Contents
......
```

### When No Files Are Found
```markdown
# No files were analyzed

Possible reasons:
- All files exceeded size limit
- No files matched include patterns
- All files matched exclude patterns
```

## Cache Directory

CodeLens caches Remote repositories in `~/.codelens/repos/` to avoid repeated downloads.

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

