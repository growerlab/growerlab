package server

import (
	"context"
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
)

type EncodedObject struct {
	ctx context.Context

	uuid     string
	repoPath string
	obj      plumbing.EncodedObject
}

func NewEncodedObject(ctx context.Context, uuid, repoPath string, obj plumbing.EncodedObject) *EncodedObject {
	return &EncodedObject{
		ctx:      ctx,
		uuid:     uuid,
		repoPath: repoPath,
		obj:      obj,
	}
}

func (o *EncodedObject) UUID() string { return o.uuid }

func (o *EncodedObject) Hash() plumbing.Hash { return o.obj.Hash() }

func (o *EncodedObject) Type() plumbing.ObjectType { return o.obj.Type() }

func (o *EncodedObject) SetType(t plumbing.ObjectType) { o.obj.SetType(t) }

func (o *EncodedObject) Size() int64 { return o.obj.Size() }

func (o *EncodedObject) SetSize(s int64) { o.obj.SetSize(s) }

func (o *EncodedObject) Reader() (io.ReadCloser, error) {
	return o.obj.Reader()
}

func (o *EncodedObject) Writer() (io.WriteCloser, error) {
	return o.obj.Writer()
}

func (o *EncodedObject) PBEncodeObject() *pb.EncodedObject {
	return &pb.EncodedObject{
		UUID: o.uuid,
		Hash: o.obj.Hash().String(),
		Type: o.obj.Type().String(),
		Size: o.obj.Size(),
	}
}
