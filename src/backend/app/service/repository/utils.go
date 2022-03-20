package repository

import "strings"

func UsernameToFilePath(owner, repoName string) string {
	p := []string{owner[:2], repoName[:2], owner, repoName + ".git"}
	return strings.Join(p, "/")
}
