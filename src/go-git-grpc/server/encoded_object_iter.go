package server

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
)

func (s *Store) NewEncodedObjectIter(ctx context.Context, tp *pb.ObjectType) (*pb.None, error) {
	var (
		uuid     = buildUUID(nil)
		result   = &pb.None{RepoPath: tp.RepoPath, UUID: uuid}
		repoPath = tp.RepoPath
	)

	err := repo(s.root, repoPath, func(repo *git.Repository) error {
		objType, err := plumbing.ParseObjectType(tp.Type)
		if err != nil {
			return errors.WithStack(err)
		}
		iter, err := repo.Storer.IterEncodedObjects(objType)
		if err != nil {
			return errors.WithStack(err)
		}

		s.putIter(&EncodedObjectIterExt{
			EncodedObjectIter: iter,
			uuid:              uuid,
		})
		return nil
	})
	return result, err
}

// EncodedObjectNext(context.Context, *None) (*EncodedObject, error)
func (s *Store) EncodedObjectNext(ctx context.Context, none *pb.None) (*pb.EncodedObject, error) {
	var (
		iterUUID = none.UUID
		repoPath = none.RepoPath
	)

	iter, ok := s.getIter(iterUUID)
	if !ok {
		return nil, ErrNotFoundIter
	}

	obj, err := iter.Next()
	if err != nil {
		return nil, err
	}

	newObj := NewEncodedObject(context.Background(), buildUUID(obj), repoPath, obj)
	s.putObject(newObj)

	return newObj.PBEncodeObject(), nil
}

func (s *Store) EncodedObjectForEach(none *pb.None, stream pb.Storer_EncodedObjectForEachServer) error {
	var (
		iterUUID = none.UUID
		repoPath = none.RepoPath
	)
	iter, ok := s.getIter(iterUUID)
	if !ok {
		return ErrNotFoundIter
	}

	err := iter.ForEach(func(object plumbing.EncodedObject) error {
		obj := NewEncodedObject(context.Background(), object.Hash().String(), repoPath, object)
		s.putObject(obj)
		return stream.Send(obj.PBEncodeObject())
	})
	return err
}

func (s *Store) EncodedObjectClose(ctx context.Context, none *pb.None) (*pb.None, error) {
	iter, ok := s.getIter(none.UUID)
	if !ok {
		return nil, ErrNotFoundIter
	}
	iter.Close()
	return none, nil
}

func (s *Store) putIter(iter *EncodedObjectIterExt) {
	s.iterStash.Set(iter)
}
func (s *Store) getIter(uuid string) (*EncodedObjectIterExt, bool) {
	iterObj, ok := s.iterStash.Get(uuid)
	if !ok {
		return nil, false
	}
	iter, ok := iterObj.(*EncodedObjectIterExt)
	return iter, ok
}
