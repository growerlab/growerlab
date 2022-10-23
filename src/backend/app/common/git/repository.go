package git

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/path"
	"github.com/growerlab/growerlab/src/go-git-grpc/client"
)

type gprcCallbackFunc func(client *client.Store) error
type gitCallbackFunc func(repo *git.Repository) error

func NewRepository(ctx context.Context, pathGroup string) *Repository {
	return &Repository{ctx: ctx, pathGroup: pathGroup}
}

type Repository struct {
	ctx       context.Context
	pathGroup string
}

func (r *Repository) CreateRepositry() error {
	r.getGit(r.ctx, r.pathGroup, func(repo *git.Repository) error {
		// TODO
		return nil
	})
	return nil
}

func (r *Repository) getGit(ctx context.Context, pathGroup string, gitcb gitCallbackFunc) error {
	err := r.getGrpcClient(ctx, pathGroup, func(client *client.Store) error {
		repo, err := git.Open(client, nil)
		if err != nil {
			return errors.Trace(err)
		}
		err = gitcb(repo)
		if err != nil {
			return errors.Trace(err)
		}
		return nil
	})
	return errors.Trace(err)
}

func (r *Repository) getGrpcClient(ctx context.Context, pathGroup string, cb gprcCallbackFunc) error {
	relativelyPath := path.GetRelativeRepositryPath(pathGroup)
	store, closeFn, err := GetGitGRPCClient(ctx, relativelyPath)
	if err != nil {
		return errors.Trace(err)
	}
	defer closeFn.Close()
	err = cb(store)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}
