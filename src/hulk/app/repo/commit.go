package repo

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	MessageMaxLen = 512 // message的最大长度
)

type Signature struct {
	// Name represents a person name. It is an arbitrary string.
	Name string `json:"name"`
	// Email is an email, but it cannot be assumed to be well-formed.
	Email string `json:"email"`
	// When is the timestamp of the signature.
	When int64 `json:"when"`
}

type PlainCommit struct {
	Hash      string    `json:"hash"`
	Author    Signature `json:"author"`
	Committer Signature `json:"committer"`
	Message   string    `json:"message"`
}

type PlainCommits []*PlainCommit

func (p PlainCommits) ToString() (string, error) {
	s, err := json.Marshal(p)
	return string(s), errors.WithStack(err)
}

func BuildPlainCommits(commits ...*object.Commit) PlainCommits {
	length := len(commits)
	if length == 0 {
		return nil
	}

	result := make([]*PlainCommit, 0, length)
	for _, cmt := range commits {
		msg := []rune(cmt.Message)
		if len(msg) > MessageMaxLen {
			msg = msg[:MessageMaxLen]
		}

		result = append(result, &PlainCommit{
			Hash: cmt.Hash.String(),
			Author: Signature{
				Name:  cmt.Author.Name,
				Email: cmt.Author.Email,
				When:  cmt.Author.When.Unix(),
			},
			Committer: Signature{
				Name:  cmt.Committer.Name,
				Email: cmt.Committer.Email,
				When:  cmt.Committer.When.Unix(),
			},
			Message: string(msg),
		})
	}
	return result
}
