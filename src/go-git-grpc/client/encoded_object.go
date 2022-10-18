package client

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
)

var _ plumbing.EncodedObject = (*EncodedObject)(nil)

type EncodedObject struct {
	ctx    context.Context
	client pb.StorerClient

	repoPath string
	uuid     string

	encodedObject *FixableEncodedObject
}

func (e *EncodedObject) Hash() plumbing.Hash {
	if e.encodedObject != nil {
		return e.encodedObject.Hash()
	}

	params := &pb.None{RepoPath: e.repoPath, UUID: e.uuid}
	resp, err := e.client.EncodedObjectHash(e.ctx, params)
	if err != nil {
		log.Printf("call EncodedObjectHash was err: %+v\n", err)
		return plumbing.ZeroHash
	}
	return plumbing.NewHash(resp.Value)
}

var UnknownObjectType plumbing.ObjectType = -126

func (e *EncodedObject) Type() plumbing.ObjectType {
	if e.encodedObject != nil {
		return e.encodedObject.Type()
	}

	params := &pb.None{RepoPath: e.repoPath, UUID: e.uuid}
	resp, err := e.client.EncodedObjectType(e.ctx, params)
	if err != nil {
		log.Printf("call EncodedObjectType was err: %+v\n", err)
		return UnknownObjectType
	}
	return plumbing.ObjectType(int8(resp.Value))
}

func (e *EncodedObject) SetType(objectType plumbing.ObjectType) {
	if e.encodedObject != nil {
		e.encodedObject.SetType(objectType)
	}

	params := &pb.Int{RepoPath: e.repoPath, UUID: e.uuid, Value: int32(objectType)}
	_, err := e.client.SetEncodedObjectType(e.ctx, params)
	if err != nil {
		log.Printf("call SetEncodedObjectType was err: %+v\n", err)
		return
	}
}

func (e *EncodedObject) Size() int64 {
	if e.encodedObject != nil {
		return e.encodedObject.Size()
	}

	params := &pb.None{RepoPath: e.repoPath, UUID: e.uuid}
	resp, err := e.client.EncodedObjectSize(e.ctx, params)
	if err != nil {
		log.Printf("call EncodedObjectSize was err: %+v\n", err)
		return 0
	}
	return resp.Value
}

func (e *EncodedObject) SetSize(i int64) {
	if e.encodedObject != nil {
		e.encodedObject.SetSize(i)
	}

	params := &pb.Int64{
		RepoPath: e.repoPath,
		UUID:     e.uuid,
		Value:    i,
	}
	_, err := e.client.SetEncodedObjectSetSize(e.ctx, params)
	if err != nil {
		log.Printf("call SetEncodedObjectSetSize was err: %+v\n", err)
		return
	}
}

func (e *EncodedObject) Reader() (io.ReadCloser, error) {
	return e.buildStream(pb.RWStream_READ)
}

func (e *EncodedObject) Writer() (io.WriteCloser, error) {
	return e.buildStream(pb.RWStream_WRITE)
}

func (e *EncodedObject) buildStream(flag pb.RWStream_FlagEnum) (*EncodedObjectStream, error) {
	stream, err := e.client.EncodedObjectRWStream(e.ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 告知服务器操作
	err = stream.Send(&pb.RWStream{
		UUID:     e.uuid,
		RepoPath: e.repoPath,
		Flag:     flag,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &EncodedObjectStream{
		uuid:     e.uuid,
		repoPath: e.repoPath,
		ctx:      e.ctx,
		stream:   stream,
	}, nil
}

var _ io.ReadWriteCloser = (*EncodedObjectStream)(nil)

type EncodedObjectStream struct {
	uuid     string
	repoPath string
	ctx      context.Context
	stream   pb.Storer_EncodedObjectRWStreamClient
}

func (e *EncodedObjectStream) Read(p []byte) (n int, err error) {
	if len(p) < bytes.MinRead {
		return 0, errors.Errorf("'p' length lt %d bytes", bytes.MinRead)
	}

	var stream *pb.RWStream
	stream, err = e.stream.Recv()
	if err != nil {
		return 0, err
	}
	n = copy(p, stream.Value)
	return
}

func (e *EncodedObjectStream) Write(p []byte) (n int, err error) {
	rw := &pb.RWStream{
		Value: p,
	}
	err = e.stream.Send(rw)
	return len(p), errors.WithStack(err)
}

func (e *EncodedObjectStream) Close() error {
	err := e.stream.CloseSend()
	return errors.WithStack(err)
}
