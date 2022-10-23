package gggrpc

import (
	"context"
	"io"

	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/client"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
	"google.golang.org/grpc"
)

func NewStoreClient(ctx context.Context, grpcServerAddr string, repoRelativePath string) (*client.Store, io.Closer, error) {
	conn, err := grpc.DialContext(ctx,
		grpcServerAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		err := errors.WithStack(err)
		return nil, nil, err
	}

	c := pb.NewStorerClient(conn)
	s := client.NewStore(ctx, repoRelativePath, c)

	return s, conn, nil
}

func NewDoorClient(ctx context.Context, grpcServerAddr string) (*client.Door, io.Closer, error) {
	conn, err := grpc.DialContext(ctx,
		grpcServerAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		err := errors.WithStack(err)
		return nil, nil, err
	}

	c := pb.NewDoorClient(conn)
	door := client.NewDoor(ctx, c)
	return door, conn, nil
}
