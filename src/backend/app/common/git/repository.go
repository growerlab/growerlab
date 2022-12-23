package git

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/growerlab/growerlab/src/backend/app/utils/logger"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/path"
	"github.com/growerlab/growerlab/src/go-git-grpc/client"
	"github.com/growerlab/growerlab/src/go-git-grpc/server/command"
)

type grpcStoreCallbackFunc func(client *client.Store) error
type grpcDoorCallbackFunc func(client *client.Door) error
type gitCallbackFunc func(repo *git.Repository) error

// var (
// 	ErrRepositoryExists = errors.New("ErrRepositoryExists")
// )

func New(ctx context.Context, pathGroup string) *Repository {
	return &Repository{
		ctx:              ctx,
		cfg:              configurator.GetConf(),
		pathGroup:        pathGroup,
		repoRelativePath: path.GetRelativeRepositoryPath(pathGroup),
		repoAbsPath:      path.GetRealRepositoryPath(pathGroup),
	}
}

type Repository struct {
	ctx              context.Context
	cfg              *configurator.Config
	pathGroup        string // abcdef/abcdef
	repoRelativePath string // ab/ab/abcdef/abcdef
	repoAbsPath      string
}

func (r *Repository) Create() error {
	if !path.CheckRepoAbsPathIsEffective(r.repoAbsPath) {
		return errors.Errorf("invalid repo path: %s, not in root dir", r.pathGroup)
	}

	if exists := r.Exists(); exists {
		return errors.AlreadyExistsError(errors.Repository, errors.AlreadyExists)
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
	logger.Info("stat: '%s' %s", r.repoAbsPath, out.String())
	return err == nil
}

// AddFile 添加文件到相应的ref。
// 通过 git 命令完成添加，主要是为了触发hook，否直接通过go生成commit对象也是可以的
func (r *Repository) AddFile(file *client.File) (string, error) {
	var commitHash string
	var err error

	if err = file.Verify(); err != nil {
		return "", errors.Trace(err)
	}

	err = r.getGrpcDoorClient(context.TODO(), func(c *client.Door) error {
		commitHash, err = c.AddOrUpdateFile(&command.Context{
			Bin:      r.cfg.GitBinPath,
			RepoPath: r.repoRelativePath,
		}, file)
		return errors.Trace(err)
	})
	return commitHash, errors.Trace(err)
}

func (r *Repository) Files(ref, dir string) ([]*FileEntity, error) {
	var result = make([]*FileEntity, 0)

	err := r.getRepo(r.ctx, r.pathGroup, func(repo *git.Repository) error {
		refer, err := repo.Reference(plumbing.NewBranchReferenceName(ref), false)
		if err != nil {
			return errors.Trace(err)
		}

		// take commit
		commit, err := repo.CommitObject(refer.Hash())
		if err != nil {
			panic(err)
		}

		tree, err := commit.Tree()
		if err != nil {
			return errors.Trace(err)
		}

		err = tree.Files().ForEach(func(file *object.File) error {
			result = append(result, buildFileEntity(file))
			return nil
		})
		return errors.Trace(err)
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	// TODO sort
	return result, nil
}

func (r *Repository) runCommand(cmd string, args []string, in io.Reader, out io.Writer) error {
	err := r.getGrpcDoorClient(r.ctx, func(client *client.Door) error {
		err := client.RunCommand(&command.Context{
			Bin:      cmd,
			Args:     args,
			In:       in,
			Out:      out,
			RepoPath: "",
			Deadline: 10 * time.Second,
		})
		return errors.Trace(err)
	})
	return errors.Trace(err)
}

func (r *Repository) getRepo(ctx context.Context, pathGroup string, gitcb gitCallbackFunc) error {
	err := r.getGrpcStoreClient(ctx, pathGroup, func(c *client.Store) error {
		repo, err := git.Open(c, nil)
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

func (r *Repository) getGrpcStoreClient(ctx context.Context, pathGroup string, cb grpcStoreCallbackFunc) error {
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

func (r *Repository) getGrpcDoorClient(ctx context.Context, cb grpcDoorCallbackFunc) error {
	store, closeFn, err := GetGitGRPCDoorClient(ctx)
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
