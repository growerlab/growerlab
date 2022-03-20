package repo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/pkg/errors"
)

const (
	ZeroRef = "0000000000000000000000000000000000000000"
)

type Repository struct {
	repoPath string
	repo     *git.Repository
}

func NewRepository(repoPath string) *Repository {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		panic(errors.WithStack(err))
	}

	return &Repository{
		repoPath: repoPath,
		repo:     repo,
	}
}

func (r *Repository) HashType(hash plumbing.Hash) plumbing.ObjectType {
	if _, err := r.repo.CommitObject(hash); err == nil {
		return plumbing.CommitObject
	}
	if _, err := r.repo.TagObject(hash); err == nil {
		return plumbing.TagObject
	}
	if _, err := r.repo.BlobObject(hash); err == nil {
		return plumbing.BlobObject
	}
	if _, err := r.repo.TreeObject(hash); err == nil {
		return plumbing.TreeObject
	}
	return plumbing.InvalidObject
}

// BetweenCommits
// maxLimit = 0 then no limit
func (r *Repository) BetweenCommits(before, after string, maxLimit uint) ([]*object.Commit, error) {
	beforeHash := plumbing.NewHash(before)
	afterHash := plumbing.NewHash(after)
	result := make([]*object.Commit, 0)

	beforeCommit, err := r.repo.CommitObject(beforeHash)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result = append(result, beforeCommit)

	count := uint(0)
	_ = beforeCommit.Parents().ForEach(func(child *object.Commit) error {
		count++
		if maxLimit > 0 && count >= maxLimit {
			return storer.ErrStop
		}
		result = append(result, child)
		if child.Hash == afterHash {
			return storer.ErrStop
		}
		return nil
	})
	return result, nil
}

func (r *Repository) TagByHash(tagHash string) (tag *object.Tag, err error) {
	tag, err = r.repo.TagObject(plumbing.NewHash(tagHash))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tag, nil
}

func (r *Repository) BranchByRef(ref string) (branch *plumbing.Reference, err error) {
	branchIter, err := r.repo.Branches()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_ = branchIter.ForEach(func(reference *plumbing.Reference) error {
		if reference.Name().String() == ref {
			branch = reference
			return storer.ErrStop
		}
		return nil
	})
	if branch == nil {
		return nil, errors.WithStack(ErrNotFoundBranch)
	}
	return
}
