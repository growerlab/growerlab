package git

import (
	"strings"

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
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Mode              filemode.FileMode `json:"-"`
	IsFile            bool              `json:"is_file"`
	TreeHash          string            `json:"tree_hash"`
	LastCommitMessage string            `json:"last_commit_message"`
	LastCommitHash    string            `json:"last_commit_hash"`
	LastCommitDate    int64             `json:"last_commit_date"`
}

func buildFileEntity(fh fileHash, commit *object.Commit) *FileEntity {
	line := strings.Split(commit.Message, "\n")
	return &FileEntity{
		ID:                fh.name,
		Name:              fh.name,
		Mode:              fh.mode,
		IsFile:            fh.mode.IsFile(),
		TreeHash:          fh.hash.String(),
		LastCommitMessage: line[0],
		LastCommitHash:    commit.Hash.String(),
		LastCommitDate:    commit.Committer.When.Unix(),
	}
}
