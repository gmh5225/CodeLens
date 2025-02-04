package git

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// SupportedGitHosts contains all supported Git hosting services
var SupportedGitHosts = []string{
	"https://github.com/",
	"https://gitlab.com/",
	"https://gitea.com/",
	"https://gitee.com/",
	"https://bitbucket.org/",
}

// IsValidGitHost checks if the given URL is from a supported Git host
func IsValidGitHost(url string) bool {
	for _, host := range SupportedGitHosts {
		if strings.HasPrefix(url, host) {
			return true
		}
	}
	return false
}

// ParseGitHubURL extracts owner and repository name from a GitHub URL
func ParseGitHubURL(url string) (owner, repo string) {
	// Remove .git suffix
	url = strings.TrimSuffix(url, ".git")

	// Remove protocol prefix
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	// Remove github.com/
	url = strings.TrimPrefix(url, "github.com/")

	// Split path
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}

	return "", ""
}

type LanguageStats struct {
	Language string
	Bytes    int
}

// GetGitHubLanguages retrieves language statistics from GitHub API
func GetGitHubLanguages(repo string) ([]LanguageStats, error) {
	owner, repoName := ParseGitHubURL(repo)
	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("invalid GitHub URL")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", owner, repoName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var langMap map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&langMap); err != nil {
		return nil, err
	}

	var stats []LanguageStats
	for lang, bytes := range langMap {
		stats = append(stats, LanguageStats{
			Language: lang,
			Bytes:    bytes,
		})
	}

	// Sort by bytes in descending order
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Bytes > stats[j].Bytes
	})

	return stats, nil
}

// GetRepoLanguages retrieves repository language statistics (currently GitHub only)
func GetRepoLanguages(repo string) ([]LanguageStats, error) {
	// Check if it's a GitHub repository
	if strings.Contains(repo, "github.com") {
		return GetGitHubLanguages(repo)
	}
	return nil, fmt.Errorf("language statistics are currently only supported for GitHub repositories")
}
