package git

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/logger"
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

func (r *Repository) TreeFiles(ref, dir string) ([]*FileEntity, error) {
	var result = make([]*FileEntity, 0)

	if dir == "/" {
		dir = ""
	}
	if strings.HasPrefix(dir, "/") {
		dir = strings.TrimPrefix(dir, "/")
	}

	err := r.getRepo(func(repo *git.Repository) error {
		var refer *plumbing.Reference
		var err error

		if plumbing.IsHash(ref) {
			refer = plumbing.NewHashReference(plumbing.ReferenceName(RefCommit), plumbing.NewHash(ref))
		} else {
			commitHash, err := repo.ResolveRevision(plumbing.Revision(ref))
			if err != nil {
				return errors.Trace(err)
			}
			refer = plumbing.NewHashReference(plumbing.ReferenceName(RefCommit), *commitHash)
		}
		if err != nil {
			return errors.Trace(err)
		}
		if refer == nil {
			return errors.NotFoundError(errors.Reference)
		}

		// not init
		if refer.Hash().IsZero() {
			return errors.RepositoryError(errors.Empty)
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

		if dir != "" {
			tree, err = tree.Tree(dir)
			if err != nil {
				return errors.Trace(err)
			}
		}

		var paths = make([]string, 0, len(tree.Entries))
		for _, entry := range tree.Entries {
			paths = append(paths, entry.Name)
		}

		nameCommitSet, err := r.getCommitForPaths(repo, commit.Hash, dir, paths)
		if err != nil {
			return errors.Trace(err)
		}
		for fileHash, cmt := range nameCommitSet {
			result = append(result, buildFileEntity(fileHash, cmt))
		}

		return errors.Trace(err)
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	// TODO sort
	return result, nil
}

func (r *Repository) listRefs(repo *git.Repository) ([]*plumbing.Reference, error) {
	var references = make([]*plumbing.Reference, 0)
	iter, err := repo.References()
	if err != nil {
		return nil, errors.Trace(err)
	}
	err = iter.ForEach(func(reference *plumbing.Reference) error {
		references = append(references, reference)
		return nil
	})
	return references, errors.Trace(err)
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

func (r *Repository) getRepo(gitcb gitCallbackFunc) error {
	err := r.getGrpcStoreClient(r.ctx, r.pathGroup, func(c *client.Store) error {
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
