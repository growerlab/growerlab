package client

import (
	"context"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/common"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
	"google.golang.org/grpc"
)

var _ storage.Storer = (*Store)(nil)

func NewStore(ctx context.Context, repoRelativePath string, pbClient pb.StorerClient) *Store {
	return &Store{
		repoPath: repoRelativePath,
		lastErr:  nil,
		ctx:      ctx,
		client:   pbClient,
	}
}

type Store struct {
	repoPath string
	lastErr  error

	ctx      context.Context
	grpcConn *grpc.ClientConn
	client   pb.StorerClient
}

func (s *Store) NewEncodedObject() plumbing.EncodedObject {
	params := &pb.None{RepoPath: s.repoPath}
	resp, err := s.client.NewEncodedObject(s.ctx, params)
	if err != nil {
		s.lastErr = errors.WithStack(err)
		return nil
	}

	return &EncodedObject{
		ctx:      s.ctx,
		client:   s.client,
		repoPath: s.repoPath,
		uuid:     resp.Value,
	}
}

func (s *Store) SetEncodedObject(obj plumbing.EncodedObject) (plumbing.Hash, error) {
	ob := obj.(*EncodedObject)
	params := &pb.UUID{Value: ob.uuid}

	hash, err := s.client.SetEncodedObject(s.ctx, params)
	if err != nil {
		return plumbing.ZeroHash, errors.WithStack(err)
	}
	return plumbing.NewHash(hash.Value), nil
}

func (s *Store) EncodedObject(objectType plumbing.ObjectType, hash plumbing.Hash) (plumbing.EncodedObject, error) {
	params := &pb.GetEncodeObject{
		RepoPath: s.repoPath,
		Hash:     hash.String(),
		Type:     objectType.String(),
	}
	obj, err := s.client.EncodedObjectEntity(s.ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), plumbing.ErrObjectNotFound.Error()) {
			return nil, plumbing.ErrObjectNotFound
		} else if strings.Contains(err.Error(), plumbing.ErrInvalidType.Error()) {
			return nil, plumbing.ErrInvalidType
		}
		return nil, errors.WithStack(err)
	}

	result := buildEncodedObjectFromPB(s.ctx, s.client, s.repoPath, obj)
	return result, nil
}

func (s *Store) IterEncodedObjects(objectType plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	iter, err := NewEncodedObjectIter(s.ctx, s.client, s.repoPath, objectType)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (s *Store) HasEncodedObject(hash plumbing.Hash) error {
	panic("implement me")
}

func (s *Store) EncodedObjectSize(hash plumbing.Hash) (int64, error) {
	params := &pb.None{
		RepoPath: s.repoPath,
		UUID:     hash.String(),
	}
	n, err := s.client.EncodedObjectSize(s.ctx, params)
	return n.Value, err
}

func (s *Store) SetReference(reference *plumbing.Reference) error {
	params := &pb.Reference{
		RepoPath: s.repoPath,
		T:        reference.Type().String(),
		N:        reference.Name().String(),
		H:        reference.Hash().String(),
		Target:   reference.Target().String(),
	}
	_, err := s.client.SetReference(s.ctx, params)
	return err
}

func (s *Store) CheckAndSetReference(new, old *plumbing.Reference) error {
	params := &pb.SetReferenceParams{
		RepoPath: s.repoPath,
		New: &pb.Reference{
			T:      new.Type().String(),
			N:      new.Name().String(),
			H:      new.Hash().String(),
			Target: new.Target().String(),
		},
		Old: &pb.Reference{
			T:      old.Type().String(),
			N:      old.Name().String(),
			H:      old.Hash().String(),
			Target: old.Target().String(),
		},
	}
	_, err := s.client.CheckAndSetReference(s.ctx, params)
	return err
}

func (s *Store) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	params := &pb.ReferenceName{
		RepoPath: s.repoPath,
		Name:     name.String(),
	}
	result, err := s.client.GetReference(s.ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), plumbing.ErrReferenceNotFound.Error()) {
			return nil, plumbing.ErrReferenceNotFound
		}
		return nil, errors.WithStack(err)
	}
	if len(result.Target) > 0 {
		return plumbing.NewReferenceFromStrings(result.N, result.Target), nil
	} else {
		return plumbing.NewHashReference(plumbing.ReferenceName(result.N), plumbing.NewHash(result.H)), nil
	}
}

const symrefPrefix = "ref: "

func (s *Store) IterReferences() (storer.ReferenceIter, error) {
	params := &pb.None{
		RepoPath: s.repoPath,
	}
	pbRefs, err := s.client.GetReferences(s.ctx, params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	refs := make([]*plumbing.Reference, 0, len(pbRefs.Refs))

	for _, r := range pbRefs.Refs {
		if r.T == plumbing.SymbolicReference.String() {
			target := symrefPrefix + r.Target
			ref := plumbing.NewReferenceFromStrings(r.N, target)
			refs = append(refs, ref)
		} else {
			ref := plumbing.NewHashReference(plumbing.ReferenceName(r.N), plumbing.NewHash(r.H))
			refs = append(refs, ref)
		}
	}

	return storer.NewReferenceSliceIter(refs), nil
}

func (s *Store) RemoveReference(name plumbing.ReferenceName) error {
	params := &pb.ReferenceName{
		RepoPath: s.repoPath,
		Name:     string(name),
	}
	_, err := s.client.RemoveReference(s.ctx, params)
	return errors.WithStack(err)
}

func (s *Store) CountLooseRefs() (int, error) {
	params := &pb.None{
		RepoPath: s.repoPath,
	}
	n, err := s.client.CountLooseRefs(s.ctx, params)
	return int(n.Value), errors.WithStack(err)
}

func (s *Store) PackRefs() error {
	params := &pb.None{
		RepoPath: s.repoPath,
	}
	_, err := s.client.PackRefs(s.ctx, params)
	return errors.WithStack(err)
}

func (s *Store) SetShallow(hashes []plumbing.Hash) error {
	if len(hashes) == 0 {
		return nil
	}
	hashesStrs := make([]string, 0, len(hashes))
	for _, h := range hashes {
		hashesStrs = append(hashesStrs, h.String())
	}

	params := &pb.Hashs{
		RepoPath: s.repoPath,
		Hash:     hashesStrs,
	}
	_, err := s.client.SetShallow(s.ctx, params)
	return errors.WithStack(err)
}

func (s *Store) Shallow() ([]plumbing.Hash, error) {
	params := &pb.None{
		RepoPath: s.repoPath,
	}
	rawHashes, err := s.client.Shallow(s.ctx, params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make([]plumbing.Hash, 0, len(rawHashes.Hash))
	for _, h := range rawHashes.Hash {
		result = append(result, plumbing.NewHash(h))
	}
	return result, nil
}

func (s *Store) SetIndex(index *index.Index) error {
	params := &pb.Index{
		RepoPath: s.repoPath,
	}
	_, err := s.client.SetIndex(s.ctx, params)
	return errors.WithStack(err)
}

func (s *Store) Index() (*index.Index, error) {
	params := &pb.None{
		RepoPath: s.repoPath,
	}
	idx, err := s.client.GetIndex(s.ctx, params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return common.BuildPbRefToIndex(idx), nil
}

func buildEncodedObjectFromPB(ctx context.Context, client pb.StorerClient, repoPath string, obj *pb.EncodedObject) plumbing.EncodedObject {
	typ, _ := plumbing.ParseObjectType(obj.Type)

	return &EncodedObject{
		ctx:      ctx,
		client:   client,
		repoPath: repoPath,
		uuid:     obj.UUID,
		encodedObject: &FixableEncodedObject{
			hash: plumbing.NewHash(obj.Hash),
			typ:  typ,
			size: obj.Size,
		},
	}
}
