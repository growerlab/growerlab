package server

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func repo(root, path string, repoFn func(*git.Repository) error) error {
	dir := filepath.Join(root, path)
	r, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}
	return repoFn(r)
}
