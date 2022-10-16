package app

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/growerlab/growerlab/src/common/configurator"

	"github.com/gliderlabs/ssh"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/mensa/app/common"
)

func NewGitSSHServer(cfg *configurator.Mensa) *GitSSHServer {
	deadline := DefaultDeadline * time.Second
	idleTimeout := DefaultIdleTimeout * time.Second

	if cfg.Deadline > 0 {
		deadline = time.Duration(cfg.Deadline) * time.Second
	}
	if cfg.IdleTimeout > 0 {
		idleTimeout = time.Duration(cfg.IdleTimeout) * time.Second
	}

	gitServer := &GitSSHServer{
		gitUser:     cfg.User,
		listen:      cfg.SSHListen,
		gitBinPath:  cfg.GitPath,
		deadline:    deadline,
		idleTimeout: idleTimeout,
	}
	return gitServer
}

type GitSSHServer struct {
	handler MiddlewareHandler

	srv         *ssh.Server
	gitBinPath  string        // bin git
	gitUser     string        // default "git"
	listen      string        // listen addr
	deadline    time.Duration // default 3600
	idleTimeout time.Duration // default 120
}

// Shutdown close all server and wait.
func (g *GitSSHServer) Shutdown() error {
	var err error
	if g.srv != nil {
		err = g.srv.Close()
	}
	return errors.WithStack(err)
}

// Start server
func (g *GitSSHServer) ListenAndServe(handler MiddlewareHandler) error {
	log.Printf("[ssh] git listen and serve: %v\n", g.listen)
	g.handler = handler

	if err := g.validate(); err != nil {
		return err
	}
	if err := g.run(); err != nil {
		return err
	}
	return nil
}

func (g *GitSSHServer) sessionHandler(session ssh.Session) {
	writeSession := &delayWriteSession{
		first:   bytes.NewBuffer(nil),
		session: session,
	}
	defer session.Close()
	go func() {
		r := session.Stderr()
		io.Copy(r, os.Stdout)
	}()

	ctx, err := common.BuildContextFromSSH(session)
	if err != nil {
		log.Printf("[ssh] %v\n", err)
		return
	}
	log.Println("[ssh] git handler commands: ", ctx.RawCommands, ctx.RepoDir)

	result := g.handler(ctx)
	if result.Err != nil {
		writeSession.Append(packetWrite(fmt.Sprintf("\x02%s\n", result.HttpMessage)))
		return
	}

	rpc, ok := AllowedCommandMap[ctx.RawCommands[0]]
	if !ok {
		log.Printf("[ssh] invalid rpc: %s\n", ctx.RawCommands[0])
		return
	}

	// 客户端push：输出到客户端的终端，之后这块应该要抽出来结构化
	if rpc == ReceivePack {
		writeSession.Append(packetWrite(fmt.Sprintf("\x02%s\n", BannerMessage)))
	}

	args := make([]string, 0)
	if rpc == ReceivePack {
		for _, op := range GitReceivePackOptions {
			args = append(args, op.Name, op.Args)
		}
	}
	args = append(args, rpc, ".")
	err = gitCommand(session, writeSession, ctx.RepoDir, args, ctx.Env())
	if err != nil {
		log.Printf("[ssh] git was err on running: %v\n", err)
	}

	// 当有修改仓库时，更新仓库
	if rpc == ReceivePack {
		err = updateServerInfo(ctx.RepoDir, ctx.Env())
	}
	log.Printf("[ssh] git handler result: %v\n", err)
}

func (g *GitSSHServer) validate() error {
	if g.gitUser == "" {
		return errors.New("git user is required")
	}
	if !strings.Contains(g.listen, ":") {
		return errors.New("invalid listen addr")
	}
	return nil
}

func (g *GitSSHServer) prepre() {
}

func (g *GitSSHServer) run() error {
	passwordOption := ssh.PasswordAuth(func(_ ssh.Context, _ string) bool {
		return false
	})

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		if g.gitUser != ctx.User() {
			return false
		}
		return true
	})

	defaultOption := func(server *ssh.Server) error {
		server.IdleTimeout = g.idleTimeout
		server.MaxTimeout = g.deadline
		server.Version = UA
		return nil
	}

	g.srv = &ssh.Server{
		Handler: g.sessionHandler,
		Addr:    g.listen,
	}
	g.srv.SetOption(publicKeyOption)
	g.srv.SetOption(passwordOption)
	g.srv.SetOption(defaultOption)
	g.srv.SetOption(ssh.NoPty())
	err := g.srv.ListenAndServe()
	return errors.WithStack(err)
}

// delayWriteSession 在unpack指令返回之前插入自定义数据
type delayWriteSession struct {
	delay   bool
	first   *bytes.Buffer
	session io.Writer
}

func (d *delayWriteSession) Append(s []byte) {
	d.first.Write(s)
}

func (d *delayWriteSession) Write(p []byte) (int, error) {
	if !d.delay && bytes.Index(p, []byte("000eunpack ok")) >= 0 {
		_, err := d.session.Write(d.first.Bytes())
		if err != nil {
			return 0, err
		}
		d.delay = true
	}
	return d.session.Write(p)
}
