package git

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type FileEntity struct {
	Name     string
	Mode     filemode.FileMode
	BlobHash plumbing.Hash
}

func buildFileEntity(file *object.File) *FileEntity {
	return &FileEntity{
		Name:     file.Name,
		Mode:     file.Mode,
		BlobHash: file.Hash,
	}
}
