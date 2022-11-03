package git

import (
	"bytes"
	"context"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/growerlab/growerlab/src/backend/app/utils/logger"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/path"
	"github.com/growerlab/growerlab/src/go-git-grpc/client"
	ggit "github.com/growerlab/growerlab/src/go-git-grpc/server/git"
)

type grpcCallbackFunc func(client *client.Store) error
type gitCallbackFunc func(repo *git.Repository) error

func New(ctx context.Context, pathGroup string) *Repository {
	return &Repository{ctx: ctx, pathGroup: pathGroup, cfg: configurator.GetConf()}
}

type Repository struct {
	ctx       context.Context
	cfg       *configurator.Config
	pathGroup string
}

func (r *Repository) CreateRepository() error {
	repoPath := path.GetRealRepositoryPath(r.pathGroup)
	door, closeFn, err := GetGitGRPCDoorClient(r.ctx)
	if err != nil {
		return errors.Trace(err)
	}
	defer closeFn.Close()

	var out bytes.Buffer
	err = door.RunGit(&ggit.Context{
		GitBin:   r.cfg.GitBinPath,
		Args:     []string{"init", "--bare", repoPath},
		In:       nil,
		Out:      &out,
		RepoPath: "",
		Deadline: 10 * time.Second, // git 执行时间
	})
	logger.Info("init repository: %s, err: %+v", repoPath, err)
	return errors.Trace(err)
}

func (r *Repository) getRepo(ctx context.Context, pathGroup string, gitcb gitCallbackFunc) error {
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

func (r *Repository) getGrpcClient(ctx context.Context, pathGroup string, cb grpcCallbackFunc) error {
	relativelyPath := path.GetRelativeRepositoryPath(pathGroup)
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
