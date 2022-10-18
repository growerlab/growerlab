package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	gggrpc "github.com/growerlab/growerlab/src/go-git-grpc"
	"github.com/growerlab/growerlab/src/go-git-grpc/client"
	"github.com/growerlab/growerlab/src/go-git-grpc/server/git"
)

// 测试 git-upload-pack git-receive-pack 操作

var repoPath = "testrepo_bare"

func init() {
	gitRoot := os.Getenv("GO_GIT_GRPC_TEST_DIR")

	go func() {
		err := gggrpc.NewServer(gitRoot, "localhost:9001")
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(1 * time.Second)
}

func main() {
	door, closeFn, err := gggrpc.NewDoorClient(context.Background(), "localhost:9001")
	if err != nil {
		panic(err)
	}
	defer closeFn.Close()

	if err := testUploadCommand(door); err != nil {
		panic(err)
	}
	if err := testPush(); err != nil {
		panic(err)
	}
}

func testUploadCommand(door *client.Door) error {
	// 测试 git-upload-pack
	in := bytes.Buffer{}
	out := bytes.Buffer{}

	cmd := &git.Context{
		Env:      []string{""},
		GitBin:   "git",
		Args:     []string{"upload-pack", "--advertise-refs", "."},
		RepoPath: repoPath,
		In:       &in,
		Out:      &out,
		Deadline: 10 * time.Second,
	}

	if err := door.RunGit(cmd); err != nil {
		return err
	}

	fmt.Println(out.String())
	time.Sleep(time.Second)
	return nil
}

func testPush() error {

	return nil
}

func testPull() error {
	return nil
}
