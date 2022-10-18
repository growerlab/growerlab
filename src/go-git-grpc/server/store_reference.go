package server

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/common"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
)

func (s *Store) SetReference(ctx context.Context, reference *pb.Reference) (*pb.None, error) {
	var result = &pb.None{}
	err := repo(s.root, reference.RepoPath, func(r *git.Repository) error {
		var ref *plumbing.Reference
		if len(reference.Target) > 0 {
			ref = plumbing.NewReferenceFromStrings(reference.N, reference.Target)
		} else {
			ref = plumbing.NewHashReference(plumbing.ReferenceName(reference.N), plumbing.NewHash(reference.H))
		}

		err := r.Storer.SetReference(ref)
		return errors.WithStack(err)
	})
	return result, err
}

func (s *Store) CheckAndSetReference(ctx context.Context, params *pb.SetReferenceParams) (*pb.None, error) {
	var result = &pb.None{}
	err := repo(s.root, params.RepoPath, func(r *git.Repository) error {
		newRef := plumbing.NewReferenceFromStrings(params.New.N, params.New.Target)
		oldRef := plumbing.NewReferenceFromStrings(params.Old.N, params.Old.Target)
		err := r.Storer.CheckAndSetReference(newRef, oldRef)
		return errors.WithStack(err)
	})
	return result, err
}

func (s *Store) GetReference(ctx context.Context, name *pb.ReferenceName) (*pb.Reference, error) {
	var result *pb.Reference
	err := repo(s.root, name.RepoPath, func(r *git.Repository) error {
		ref, err := r.Storer.Reference(plumbing.ReferenceName(name.Name))
		if err != nil {
			return errors.WithStack(err)
		}

		result = common.BuildRefToPbRef(ref)
		return nil
	})
	return result, err
}

func (s *Store) GetReferences(ctx context.Context, none *pb.None) (*pb.References, error) {
	var result = new(pb.References)
	err := repo(s.root, none.RepoPath, func(r *git.Repository) error {
		iter, err := r.Storer.IterReferences()
		if err != nil {
			return errors.WithStack(err)
		}

		err = iter.ForEach(func(ref *plumbing.Reference) error {
			pbRef := common.BuildRefToPbRef(ref)
			result.Refs = append(result.Refs, pbRef)
			return nil
		})
		return errors.WithStack(err)
	})
	return result, err
}

func (s *Store) RemoveReference(ctx context.Context, name *pb.ReferenceName) (*pb.None, error) {
	err := repo(s.root, name.RepoPath, func(r *git.Repository) error {
		rn := plumbing.ReferenceName(name.Name)
		err := r.Storer.RemoveReference(rn)
		return errors.WithStack(err)
	})
	return &pb.None{}, err
}

func (s *Store) CountLooseRefs(ctx context.Context, none *pb.None) (*pb.Int64, error) {
	var result = new(pb.Int64)
	err := repo(s.root, none.RepoPath, func(r *git.Repository) error {
		count, err := r.Storer.CountLooseRefs()
		if err != nil {
			return errors.WithStack(err)
		}
		result.Value = int64(count)
		return nil
	})
	return result, err
}

func (s *Store) PackRefs(ctx context.Context, none *pb.None) (*pb.None, error) {
	var result = new(pb.None)
	err := repo(s.root, none.RepoPath, func(r *git.Repository) error {
		return errors.WithStack(r.Storer.PackRefs())
	})
	return result, err
}
