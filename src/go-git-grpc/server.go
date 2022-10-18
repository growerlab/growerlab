package gggrpc

import (
	"net"

	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
	"github.com/growerlab/growerlab/src/go-git-grpc/server"
	"google.golang.org/grpc"
)

func NewServer(root, address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	store := server.NewStore(root)
	door := server.NewDoor(root) // NOTE door 目前放在一个进程中，如果之后发现store和door会相互影响的话，可以拆分

	pb.RegisterStorerServer(s, store)
	pb.RegisterDoorServer(s, door)
	if err := s.Serve(lis); err != nil {
		return errors.Errorf("failed to serve: %v", err)
	}
	return nil
}
