package git

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type ReferenceType string

const (
	RefBranch ReferenceType = "branch"
	RefTag    ReferenceType = "tag"
	RefCommit ReferenceType = "commit"
)

type FileEntity struct {
	Name     string            `json:"name"`
	Mode     filemode.FileMode `json:"-"`
	BlobHash plumbing.Hash     `json:"blob_hash"`
}

func buildFileEntity(file *object.File) *FileEntity {
	return &FileEntity{
		Name:     file.Name,
		Mode:     file.Mode,
		BlobHash: file.Hash,
	}
}
