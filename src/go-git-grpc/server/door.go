package server

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
	"github.com/growerlab/growerlab/src/go-git-grpc/server/command"
)

// ServerCommand is used for a single server command execution.
type ServerCommand struct {
	repoPath     string
	gitBinServer pb.Door_RunCommandServer

	readBuf *bytes.Buffer
	ctx     *command.Context
}

// Start 协议：第一个请求仅传git相关的参数，不传数据
func (s *ServerCommand) Start() error {
	firstReq, err := s.gitBinServer.Recv()
	if err != nil {
		return err
	}

	var in io.Reader
	var out io.Writer

	if firstReq.HasInput {
		in = s
	}
	if firstReq.HasOutput {
		out = s
	}

	s.ctx = &command.Context{
		Env:      firstReq.Env,
		Bin:      firstReq.Bin,
		Args:     firstReq.Args,
		In:       in,
		Out:      out,
		RepoPath: firstReq.Path,
		Deadline: time.Duration(firstReq.Deadline),
	}
	s.repoPath = firstReq.Path
	s.readBuf = bytes.NewBufferString("")

	log.Println("---->", s.ctx.String())

	return nil
}

func (s *ServerCommand) Read(p []byte) (n int, err error) {
	if s.readBuf.Len() > 0 {
		return s.readBuf.Read(p)
	}
	req, err := s.gitBinServer.Recv()
	if err != nil {
		return 0, err
	}
	if len(req.Raw) == 0 {
		return 0, nil
	}
	n, err = s.readBuf.Write(req.Raw)
	if err != nil {
		return 0, err
	}
	if n > 0 {
		return s.Read(p)
	}
	return 0, nil
}

func (s *ServerCommand) Write(p []byte) (n int, err error) {
	err = s.gitBinServer.Send(&pb.Response{Raw: p})
	return len(p), err
}

func (s *ServerCommand) Close() error {
	return nil
}

var _ pb.DoorServer = (*Door)(nil)

func NewDoor(root string) *Door {
	return &Door{
		root: root,
	}
}

type Door struct {
	*pb.UnimplementedDoorServer
	root string
}

// RunCommand 执行git命令
func (d *Door) RunCommand(pack pb.Door_RunCommandServer) error {
	srvCmd := ServerCommand{gitBinServer: pack}
	if err := srvCmd.Start(); err != nil {
		return err
	}

	return command.Run(d.root, srvCmd.ctx)
}

func (d *Door) mustEmbedUnimplementedDoorServer() {}
