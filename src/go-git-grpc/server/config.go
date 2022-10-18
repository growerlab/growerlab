package server

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/common"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
)

func (s *Store) GetConfig(ctx context.Context, none *pb.None) (*pb.Config, error) {
	var result = new(pb.Config)
	err := repo(s.root, none.RepoPath, func(r *git.Repository) error {
		cfg, err := r.Storer.Config()
		if err != nil {
			return errors.WithStack(err)
		}
		result = common.BuildConfigFromPbConfig(cfg)
		return nil
	})
	return result, err
}

func (s *Store) SetConfig(ctx context.Context, c *pb.Config) (*pb.None, error) {
	var result = new(pb.None)
	err := repo(s.root, c.RepoPath, func(r *git.Repository) error {
		cfg := common.BuildPbConfigFromConfig(c)
		err := r.Storer.SetConfig(cfg)
		return errors.WithStack(err)
	})
	return result, err
}

func (s *Store) Modules(ctx context.Context, none *pb.None) (*pb.ModuleNames, error) {
	var result = new(pb.ModuleNames)
	err := repo(s.root, none.RepoPath, func(r *git.Repository) error {
		cfg, err := r.Storer.Config()
		if err != nil {
			return errors.WithStack(err)
		}
		for _, submd := range cfg.Submodules {
			result.Names = append(result.Names, submd.Name)
		}
		return nil
	})
	return result, err
}
