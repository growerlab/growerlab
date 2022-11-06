package git

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/growerlab/growerlab/src/backend/app/utils/logger"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/path"
	"github.com/growerlab/growerlab/src/go-git-grpc/client"
	"github.com/growerlab/growerlab/src/go-git-grpc/server/command"
)

type grpcCallbackFunc func(client *client.Store) error
type gitCallbackFunc func(repo *git.Repository) error

var (
	ErrRepositoryExists = errors.New("ErrRepositoryExists")
)

func New(ctx context.Context, pathGroup string) *Repository {
	return &Repository{
		ctx:         ctx,
		cfg:         configurator.GetConf(),
		pathGroup:   pathGroup,
		repoAbsPath: path.GetRealRepositoryPath(pathGroup),
	}
}

type Repository struct {
	ctx         context.Context
	cfg         *configurator.Config
	pathGroup   string
	repoAbsPath string
}

func (r *Repository) Create() error {
	if exists := r.Exists(); exists {
		return ErrRepositoryExists
	}

	var out bytes.Buffer
	err := r.runCommand(r.cfg.GitBinPath, []string{"init", "--bare", r.repoAbsPath}, nil, &out)
	logger.Info("init repository: '%s', git result: '%s', err: %+v", r.pathGroup, out.String(), err)
	return errors.Trace(err)
}

// Delete 删除仓库
func (r *Repository) Delete() error {
	if !path.CheckRepoAbsPathIsEffective(r.repoAbsPath) {
		return errors.Errorf("invalid repo path: %s", r.pathGroup)
	}

	err := r.runCommand("rm", []string{"-rf", r.repoAbsPath}, nil, nil)
	logger.Info("delete repository: '%s', err: %+v", r.pathGroup, err)
	return errors.Trace(err)
}

func (r *Repository) Exists() bool {
	var out bytes.Buffer
	err := r.runCommand("stat", []string{r.repoAbsPath}, nil, &out)
	if err == nil {
		return true
	}
	return false
}

func (r *Repository) runCommand(cmd string, args []string, in io.Reader, out io.Writer) error {
	door, closeFn, err := GetGitGRPCDoorClient(r.ctx)
	if err != nil {
		return errors.Trace(err)
	}
	defer closeFn.Close()

	err = door.RunCommand(&command.Context{
		Bin:      cmd,
		Args:     args,
		In:       in,
		Out:      out,
		RepoPath: "",
		Deadline: 10 * time.Second,
	})
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
