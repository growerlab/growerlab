package client

import (
	"io"

	"github.com/go-git/go-git/v5/plumbing"
)

var _ plumbing.EncodedObject = (*FixableEncodedObject)(nil)

// 仅仅用与存储基本信息的 EncodedObject 对象
// 无法进行读写（Reader、Writer）
type FixableEncodedObject struct {
	hash plumbing.Hash
	typ  plumbing.ObjectType
	size int64
}

func (t *FixableEncodedObject) Hash() plumbing.Hash { return t.hash }

func (t *FixableEncodedObject) Type() plumbing.ObjectType { return t.typ }

func (t *FixableEncodedObject) SetType(plumbing.ObjectType) {}

func (t *FixableEncodedObject) Size() int64 { return t.size }

func (t *FixableEncodedObject) SetSize(int64) {}

func (t *FixableEncodedObject) Reader() (io.ReadCloser, error) { return nil, nil }

func (t *FixableEncodedObject) Writer() (io.WriteCloser, error) { return nil, nil }
