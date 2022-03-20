package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	gggrpc "github.com/growerlab/go-git-grpc"
	"github.com/growerlab/go-git-grpc/server/git"
	"github.com/growerlab/growerlab/src/mensa/app/conf"
	"github.com/pkg/errors"
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
	gitBinPath := conf.GetConfig().GitPath
	deadline := time.Duration(conf.GetConfig().Deadline) * time.Second
	gogitgrpcAddr := conf.GetConfig().GoGitGrpcAddr

	// deadline
	cmdCtx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	gitDoor, doorWatcher, err := gggrpc.NewDoorClient(cmdCtx, gogitgrpcAddr)
	if err != nil {
		return errors.WithStack(err)
	}
	defer doorWatcher.Close()

	// for debug
	// if out != nil {
	// 	out = io.MultiWriter(os.Stdout, out)
	// }

	err = gitDoor.RunGit(&git.Context{
		Env:      envs,
		GitBin:   gitBinPath,
		Args:     args,
		In:       in,
		Out:      out,
		RepoPath: repoDir,
		Deadline: deadline,
	})
	return err
}
