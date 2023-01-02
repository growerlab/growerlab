package git

import (
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
	Name              string            `json:"name"`
	Mode              filemode.FileMode `json:"-"`
	IsFile            bool              `json:"is_file"`
	BlobHash          string            `json:"blob_hash"`
	LastCommitMessage string            `json:"last_commit_message"`
	LastCommitHash    string            `json:"last_commit_hash"`
	LastCommitDate    int64             `json:"last_commit_date"`
}

func buildFileEntity(file *object.File) *FileEntity {
	return &FileEntity{
		Name:     file.Name,
		Mode:     file.Mode,
		BlobHash: file.Hash.String(),
	}
}
