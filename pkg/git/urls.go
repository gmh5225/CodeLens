package git

import "strings"

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
