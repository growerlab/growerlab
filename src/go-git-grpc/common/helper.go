package common

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	format "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
)

func BuildRefToPbRef(ref *plumbing.Reference) *pb.Reference {
	result := &pb.Reference{
		T:      ref.Type().String(),
		N:      string(ref.Name()),
		H:      ref.Hash().String(),
		Target: string(ref.Target()),
	}
	return result
}

func BuildIndexToPbRef(idx *index.Index) *pb.Index {
	var entries []*pb.Entry
	var cache *pb.Tree
	var resolveUndo *pb.ResolveUndo
	var endOfIndexEntry *pb.EndOfIndexEntry

	if len(idx.Entries) > 0 {
		for _, ent := range idx.Entries {
			entries = append(entries, &pb.Entry{
				Hash:         ent.Hash.String(),
				Name:         ent.Name,
				CreatedAt:    ent.CreatedAt.Unix(),
				ModifiedAt:   ent.ModifiedAt.Unix(),
				Dev:          ent.Dev,
				Inode:        ent.Inode,
				Mode:         uint32(ent.Mode),
				UID:          ent.UID,
				GID:          ent.GID,
				Size:         ent.Size,
				Stage:        int64(ent.Stage),
				SkipWorktree: ent.SkipWorktree,
				IntentToAdd:  ent.IntentToAdd,
			})
		}
	}
	if idx.Cache != nil {
		cache = &pb.Tree{
			Entries: make([]*pb.TreeEntry, 0, len(idx.Cache.Entries)),
		}
		for _, ent := range idx.Cache.Entries {
			cache.Entries = append(cache.Entries, &pb.TreeEntry{
				Path:    ent.Path,
				Entries: int64(ent.Entries),
				Trees:   int64(ent.Trees),
				Hash:    ent.Hash.String(),
			})
		}
	}
	if idx.ResolveUndo != nil {
		resolveUndo = &pb.ResolveUndo{
			Entries: make([]*pb.ResolveUndoEntry, 0, len(idx.ResolveUndo.Entries)),
		}
		for _, ent := range idx.ResolveUndo.Entries {
			stages := make([]*pb.MapFieldEntry, 0, len(ent.Stages))
			for stage, hash := range ent.Stages {
				stages = append(stages, &pb.MapFieldEntry{
					Key:   int64(stage),
					Value: hash.String(),
				})
			}
			resolveUndo.Entries = append(resolveUndo.Entries, &pb.ResolveUndoEntry{
				Path:   ent.Path,
				Stages: stages,
			})
		}
	}
	if idx.EndOfIndexEntry != nil {
		endOfIndexEntry = &pb.EndOfIndexEntry{
			Offset: idx.EndOfIndexEntry.Offset,
			Hash:   idx.EndOfIndexEntry.Hash.String(),
		}
	}

	newIdx := &pb.Index{
		Version:         idx.Version,
		Entries:         entries,
		Cache:           cache,
		ResolveUndo:     resolveUndo,
		EndOfIndexEntry: endOfIndexEntry,
	}
	return newIdx
}

func BuildPbRefToIndex(idx *pb.Index) *index.Index {
	var entries []*index.Entry
	var cache *index.Tree
	var trees []index.TreeEntry
	var resolveUndo *index.ResolveUndo
	var resolveUndoEntries []index.ResolveUndoEntry
	var endOfIndexEntry *index.EndOfIndexEntry

	if len(idx.Entries) > 0 {
		entries = make([]*index.Entry, 0, len(idx.Entries))
		for _, ent := range idx.Entries {
			entries = append(entries, &index.Entry{
				Hash:         plumbing.NewHash(ent.Hash),
				Name:         ent.Name,
				CreatedAt:    time.Unix(ent.CreatedAt, 0),
				ModifiedAt:   time.Unix(ent.ModifiedAt, 0),
				Dev:          ent.Dev,
				Inode:        ent.Inode,
				Mode:         filemode.FileMode(ent.Mode),
				UID:          ent.UID,
				GID:          ent.GID,
				Size:         ent.Size,
				Stage:        index.Stage(ent.Stage),
				SkipWorktree: ent.SkipWorktree,
				IntentToAdd:  ent.IntentToAdd,
			})
		}
	}

	if idx.Cache != nil && len(idx.Cache.Entries) > 0 {
		trees = make([]index.TreeEntry, 0, len(idx.Cache.Entries))
		for _, ent := range idx.Cache.Entries {
			trees = append(trees, index.TreeEntry{
				Path:    ent.Path,
				Entries: int(ent.Entries),
				Trees:   int(ent.Trees),
				Hash:    plumbing.NewHash(ent.Hash),
			})
		}
	}
	if trees != nil {
		cache = &index.Tree{
			Entries: trees,
		}
	}

	if idx.ResolveUndo != nil {
		for _, ent := range idx.ResolveUndo.Entries {
			stageSet := make(map[index.Stage]plumbing.Hash)
			for _, stg := range ent.Stages {
				stageSet[index.Stage(stg.Key)] = plumbing.NewHash(stg.Value)
			}
			resolveUndoEntries = append(resolveUndoEntries, index.ResolveUndoEntry{
				Path:   ent.Path,
				Stages: stageSet,
			})
		}
	}
	if len(resolveUndoEntries) > 0 {
		resolveUndo = &index.ResolveUndo{Entries: resolveUndoEntries}
	}
	if idx.EndOfIndexEntry != nil {
		endOfIndexEntry = &index.EndOfIndexEntry{
			Offset: idx.EndOfIndexEntry.Offset,
			Hash:   plumbing.NewHash(idx.EndOfIndexEntry.Hash),
		}
	}

	newIdx := &index.Index{
		Version:         idx.Version,
		Entries:         entries,
		Cache:           cache,
		ResolveUndo:     resolveUndo,
		EndOfIndexEntry: endOfIndexEntry,
	}
	return newIdx
}

func BuildConfigFromPbConfig(cfg *config.Config) *pb.Config {
	var result = new(pb.Config)
	result.Core = &pb.Config_MsgCore{
		IsBare:      cfg.Core.IsBare,
		Worktree:    cfg.Core.Worktree,
		CommentChar: cfg.Core.CommentChar,
	}
	result.User = &pb.Config_MsgUser{
		Name:  cfg.User.Name,
		Email: cfg.User.Email,
	}
	result.Author = &pb.Config_MsgAuthor{
		Name:  cfg.Author.Name,
		Email: cfg.Author.Email,
	}
	result.Committer = &pb.Config_MsgCommitter{
		Name:  cfg.Committer.Name,
		Email: cfg.Committer.Email,
	}
	result.Pack = &pb.Config_MsgPack{
		Window: uint64(cfg.Pack.Window),
	}
	raw, err := json.Marshal(cfg.Raw)
	if err == nil {
		result.Raw = raw
	}

	var remotes = make([]*pb.MapRemotes, 0, len(cfg.Remotes))
	for key, remote := range cfg.Remotes {
		fetch := make([]string, 0, len(remote.Fetch))
		for i := range remote.Fetch {
			fetch = append(fetch, string(remote.Fetch[i]))
		}
		cfg := &pb.RemoteConfig{
			Name:  remote.Name,
			URLs:  remote.URLs,
			Fetch: fetch,
		}
		remotes = append(remotes, &pb.MapRemotes{
			Key:    key,
			Config: cfg,
		})
	}
	result.Remotes = remotes

	var submodules = make([]*pb.MapSubmodule, 0, len(cfg.Submodules))
	for key, submd := range cfg.Submodules {
		pbSubMD := &pb.Submodule{
			Name:   submd.Name,
			Path:   submd.Path,
			URL:    submd.URL,
			Branch: submd.Branch,
		}
		submodules = append(submodules, &pb.MapSubmodule{
			Key: key,
			Sub: pbSubMD,
		})
	}
	result.Submodules = submodules

	var branches = make([]*pb.MapBranch, 0, len(cfg.Branches))
	for key, br := range cfg.Branches {
		branch := &pb.Branch{
			Name:   br.Name,
			Remote: br.Remote,
			Merge:  string(br.Merge),
			Rebase: br.Rebase,
		}
		branches = append(branches, &pb.MapBranch{
			Key:    key,
			Branch: branch,
		})
	}
	result.Branches = branches

	return result
}

func BuildPbConfigFromConfig(cfg *pb.Config) *config.Config {
	var (
		result = &config.Config{}
		err    error
	)

	result.Raw = format.New()
	err = json.NewDecoder(bytes.NewBuffer(cfg.Raw)).Decode(result.Raw)
	if err != nil {
		panic(err)
	}

	result.Core.IsBare = cfg.Core.IsBare
	result.Core.CommentChar = cfg.Core.CommentChar
	result.Core.Worktree = cfg.Core.Worktree

	result.User.Name = cfg.User.Name
	result.User.Email = cfg.User.Email

	result.Author.Name = cfg.Author.Name
	result.Author.Email = cfg.Author.Email

	result.Committer.Name = cfg.Committer.Name
	result.Committer.Email = cfg.Committer.Email

	result.Pack.Window = uint(cfg.Pack.Window)

	result.Remotes = map[string]*config.RemoteConfig{}
	for _, remote := range cfg.Remotes {
		cfg := remote.Config
		if cfg != nil {
			fetch := make([]config.RefSpec, 0, len(cfg.Fetch))
			for _, ft := range cfg.Fetch {
				fetch = append(fetch, config.RefSpec(ft))
			}
			rmt := &config.RemoteConfig{
				Name:  cfg.Name,
				URLs:  cfg.URLs,
				Fetch: fetch,
			}
			result.Remotes[remote.Key] = rmt
		}
	}

	result.Submodules = map[string]*config.Submodule{}
	for _, submd := range cfg.Submodules {
		sb := submd.Sub
		if sb != nil {
			sub := &config.Submodule{
				Name:   sb.Name,
				Path:   sb.Path,
				URL:    sb.URL,
				Branch: sb.Branch,
			}
			result.Submodules[submd.Key] = sub
		}
	}

	result.Branches = map[string]*config.Branch{}
	for _, branch := range cfg.Branches {
		br := branch.Branch
		result.Branches[branch.Key] = &config.Branch{
			Name:   br.Name,
			Remote: br.Remote,
			Merge:  plumbing.ReferenceName(br.Merge),
			Rebase: br.Rebase,
		}
	}

	return result
}
