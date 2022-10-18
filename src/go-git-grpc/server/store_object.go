package server

import (
	"bytes"
	"context"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
)

func (s *Store) NewEncodedObject(ctx context.Context, none *pb.None) (*pb.UUID, error) {
	obj := NewEncodedObject(ctx, buildUUID(nil), none.RepoPath, &plumbing.MemoryObject{})
	s.putObject(obj)
	return &pb.UUID{Value: obj.UUID()}, nil
}

func (s *Store) SetEncodedObject(ctx context.Context, uuid *pb.UUID) (*pb.Hash, error) {
	var result *pb.Hash
	var key = uuid.Value
	var obj, exists = s.getObject(key)
	if !exists {
		return nil, ErrNotFoundObject
	}

	err := repo(s.root, obj.repoPath, func(r *git.Repository) error {
		h, err := r.Storer.SetEncodedObject(obj)
		if err != nil {
			return errors.WithStack(err)
		}
		result = &pb.Hash{
			Value: h.String(),
		}
		return nil
	})
	return result, errors.WithStack(err)
}

func (s *Store) SetEncodedObjectType(ctx context.Context, i *pb.Int) (*pb.None, error) {
	var result = &pb.None{UUID: i.UUID}
	var objectType = plumbing.ObjectType(i.Value)

	obj, exists := s.getObject(i.UUID)
	if !exists {
		return nil, ErrNotFoundObject
	}
	obj.SetType(objectType)

	return result, nil
}

func (s *Store) SetEncodedObjectSetSize(ctx context.Context, i *pb.Int64) (*pb.None, error) {
	var result = &pb.None{UUID: i.UUID}

	obj, exists := s.getObject(i.UUID)
	if !exists {
		return nil, ErrNotFoundObject
	}
	obj.SetSize(i.Value)

	return result, nil
}

func (s *Store) EncodedObjectEntity(ctx context.Context, objEntity *pb.GetEncodeObject) (*pb.EncodedObject, error) {
	var result *pb.EncodedObject
	var repoPath = objEntity.RepoPath

	err := repo(s.root, objEntity.RepoPath, func(r *git.Repository) error {
		var (
			objectType = plumbing.AnyObject
			err        error
		)
		if objEntity.Type != objectType.String() {
			objectType, err = plumbing.ParseObjectType(objEntity.Type)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		hash := plumbing.NewHash(objEntity.Hash)
		obj, err := r.Storer.EncodedObject(objectType, hash)
		if err != nil {
			return err
		}

		newObj := NewEncodedObject(ctx, buildUUID(obj), repoPath, obj)
		s.putObject(newObj)

		result = newObj.PBEncodeObject()
		return nil
	})
	return result, err
}

func (s *Store) EncodedObjectType(ctx context.Context, none *pb.None) (*pb.Int, error) {
	obj, exists := s.getObject(none.UUID)
	if !exists {
		return nil, ErrNotFoundObject
	}
	return &pb.Int{Value: int32(obj.Type())}, nil
}

func (s *Store) EncodedObjectHash(ctx context.Context, none *pb.None) (*pb.Hash, error) {
	obj, exists := s.getObject(none.UUID)
	if !exists {
		return nil, ErrNotFoundObject
	}
	return &pb.Hash{Value: obj.uuid}, nil
}

func (s *Store) EncodedObjectSize(ctx context.Context, none *pb.None) (*pb.Int64, error) {
	obj, exists := s.getObject(none.UUID)
	if !exists {
		return nil, ErrNotFoundObject
	}
	return &pb.Int64{Value: obj.Size()}, nil
}

func (s *Store) EncodedObjectRWStream(stream pb.Storer_EncodedObjectRWStreamServer) error {
	first, err := stream.Recv()
	if err != nil {
		return errors.WithStack(err)
	}
	key := first.UUID
	obj, exists := s.getObject(key)
	if !exists {
		return ErrNotFoundObject
	}

	// NOTE 这里加上同一个object对象不应该同时读写，所以不进行读写goroutine
	switch first.Flag {
	case pb.RWStream_READ:
		reader, err := obj.Reader()
		if err != nil {
			return errors.WithStack(err)
		}
		defer reader.Close()

		var buf = make([]byte, bytes.MinRead)
		for {
			var n int
			n, err = reader.Read(buf)
			if err == io.EOF {
				return nil
			}
			buf = buf[:n]
			err = stream.Send(&pb.RWStream{
				Value: buf,
			})
			if err != nil {
				return err
			}
			buf = buf[:bytes.MinRead]
		}
	case pb.RWStream_WRITE:
		writer, err := obj.Writer()
		if err != nil {
			return errors.WithStack(err)
		}
		defer writer.Close()

		for {
			rw, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			_, err = writer.Write(rw.Value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Store) getObject(uuid string) (*EncodedObject, bool) {
	obj, ok := s.objectStash.Get(uuid)
	if !ok {
		return nil, false
	}
	o, ok := obj.(*EncodedObject)
	return o, ok
}

func (s *Store) putObject(obj *EncodedObject) {
	s.objectStash.Set(obj)
}

type EncodedObjectIterExt struct {
	storer.EncodedObjectIter
	uuid string
}

func (e *EncodedObjectIterExt) UUID() string {
	return e.uuid
}
