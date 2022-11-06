package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/errors"
	gggrpc "github.com/growerlab/growerlab/src/go-git-grpc"
	"github.com/growerlab/growerlab/src/go-git-grpc/server/command"
)

type Option struct {
	Name string
	Args string
}

var GitReceivePackOptions = []*Option{
	{"-c", fmt.Sprintf("core.hooksPath=%s", filepath.Join(os.Args[0], "hooks"))},
	{"-c", "core.alternateRefsCommand=exit 0 #"},
	{"-c", "receive.fsck.badTimezone=ignore"},
}

func gitCommand(in io.Reader, out io.Writer, repoDir string, args []string, envs []string) error {
	global := configurator.GetConf()
	conf := global.Mensa
	gitBinPath := global.GitBinPath
	deadline := time.Duration(conf.Deadline) * time.Second
	goGitGrpcServerAddr := global.GoGitGrpcServerAddr

	// deadline
	cmdCtx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	gitDoor, doorWatcher, err := gggrpc.NewDoorClient(cmdCtx, goGitGrpcServerAddr)
	if err != nil {
		return errors.WithStack(err)
	}
	defer doorWatcher.Close()

	// for debug
	// if out != nil {
	// 	out = io.MultiWriter(os.Stdout, out)
	// }

	err = gitDoor.RunCommand(&command.Context{
		Env:      envs,
		Bin:      gitBinPath,
		Args:     args,
		In:       in,
		Out:      out,
		RepoPath: repoDir,
		Deadline: deadline,
	})
	return err
}
