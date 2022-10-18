package client

import (
	"encoding/json"
	"log"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	format "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/go-git/go-git/v5/storage"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/common"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
)

func (s *Store) Config() (*config.Config, error) {
	params := &pb.None{
		RepoPath: s.repoPath,
		UUID:     "",
	}
	cfg, err := s.client.GetConfig(s.ctx, params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	remotes := map[string]*config.RemoteConfig{}
	submodules := map[string]*config.Submodule{}
	branches := map[string]*config.Branch{}
	raw := new(format.Config)
	err = json.Unmarshal(cfg.Raw, raw)
	if err != nil {
		log.Printf("Config() unmarshal was err:%+v\n", err)
	}

	for _, r := range cfg.Remotes {
		var fetches = make([]config.RefSpec, 0, len(r.Config.Fetch))
		for _, f := range r.Config.Fetch {
			fetches = append(fetches, config.RefSpec(f))
		}
		remotes[r.Key] = &config.RemoteConfig{
			Name:  r.Config.Name,
			URLs:  r.Config.URLs,
			Fetch: fetches,
		}
	}

	for _, sub := range cfg.Submodules {
		s := &config.Submodule{
			Name:   sub.Sub.Name,
			Path:   sub.Sub.Path,
			URL:    sub.Sub.URL,
			Branch: sub.Sub.Branch,
		}
		submodules[sub.Key] = s
	}

	for _, bn := range cfg.Branches {
		b := &config.Branch{
			Name:   bn.Branch.Name,
			Remote: bn.Branch.Remote,
			Merge:  plumbing.NewBranchReferenceName(bn.Branch.Merge),
			Rebase: bn.Branch.Rebase,
		}
		branches[bn.Key] = b
	}

	return &config.Config{
		Core: struct {
			IsBare      bool
			Worktree    string
			CommentChar string
		}{
			IsBare:      cfg.Core.IsBare,
			Worktree:    cfg.Core.Worktree,
			CommentChar: cfg.Core.CommentChar,
		},
		User: struct {
			Name  string
			Email string
		}{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
		},
		Author: struct {
			Name  string
			Email string
		}{
			Name:  cfg.Author.Name,
			Email: cfg.Author.Email,
		},
		Committer: struct {
			Name  string
			Email string
		}{
			Name:  cfg.Committer.Name,
			Email: cfg.Committer.Email,
		},
		Pack: struct {
			Window uint
		}{
			Window: uint(cfg.Pack.Window),
		},
		Remotes:    remotes,
		Submodules: submodules,
		Branches:   branches,
		Raw:        raw,
	}, nil
}

func (s *Store) SetConfig(config *config.Config) error {
	pbConfig := common.BuildConfigFromPbConfig(config)

	params := &pb.Config{
		RepoPath:   s.repoPath,
		Remotes:    pbConfig.Remotes,
		Submodules: pbConfig.Submodules,
		Branches:   pbConfig.Branches,
		Raw:        pbConfig.Raw,
		Core:       pbConfig.Core,
		User:       pbConfig.User,
		Author:     pbConfig.Author,
		Committer:  pbConfig.Committer,
		Pack:       pbConfig.Pack,
	}
	_, err := s.client.SetConfig(s.ctx, params)
	return err
}

func (s *Store) Module(name string) (storage.Storer, error) {
	// params := &pb.None{
	// 	RepoPath: s.repoPath,
	// }
	// s.client.Modules(s.ctx, params)
	panic("未实现")
}
