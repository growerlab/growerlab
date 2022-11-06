package command

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/growerlab/growerlab/src/common/errors"
)

const DefaultTimeout = 60 * time.Minute // 推送和拉取的最大执行时间

type Context struct {
	Env      []string  // 环境变量 key=value
	Bin      string    // binary command path
	Args     []string  // args
	In       io.Reader // input
	Out      io.Writer // output
	RepoPath string    // repo dir

	Deadline time.Duration // 命令执行时间，单位秒
}

func (p *Context) String() string {
	var sb strings.Builder
	sb.WriteString("bin: ")
	sb.WriteString(p.Bin)
	sb.WriteString("\n")
	sb.WriteString("args: ")
	sb.WriteString(strings.Join(p.Args, " "))
	sb.WriteString("\n")
	sb.WriteString("repo path: ")
	sb.WriteString(p.RepoPath)
	sb.WriteString("\n")
	sb.WriteString("env: ")
	sb.WriteString(strings.Join(p.Env, " "))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("deadline: %v", p.Deadline))
	return sb.String()
}

func Run(root string, params *Context) error {
	if params.Deadline <= 0 {
		params.Deadline = DefaultTimeout
	}
	if len(params.RepoPath) > 0 {
		root = filepath.Join(root, params.RepoPath)
	}

	// deadline
	cmdCtx, cancel := context.WithTimeout(context.Background(), params.Deadline)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, params.Bin, params.Args...)
	if len(params.Env) > 0 {
		systemEnvs := os.Environ()
		cmd.Env = make([]string, 0, len(params.Env)+len(systemEnvs))
		for _, v := range params.Env {
			cmd.Env = append(cmd.Env, v)
		}
		cmd.Env = append(cmd.Env, os.Environ()...)
	}
	cmd.Dir = root
	if params.In != nil {
		cmd.Stdin = params.In
	}
	if params.Out != nil {
		cmd.Stdout = params.Out
	}

	cmd.Stderr = io.MultiWriter(log.Writer(), os.Stderr)
	err := cmd.Start()
	if err != nil {
		return errors.Wrap(err, "start git command failed")
	}

	err = cmd.Wait()
	log.Printf("command was done: %d\n", cmd.ProcessState.ExitCode())
	return errors.WithStack(err)
}
