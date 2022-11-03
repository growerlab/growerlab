package test

import (
	"os"
	"path/filepath"
	"time"

	"github.com/growerlab/growerlab/src/common"

	"github.com/growerlab/growerlab/src/common/configurator"
	gggrpc "github.com/growerlab/growerlab/src/go-git-grpc"
)

func InitForTest() {
	start(initChdir)
	start(configurator.InitConfig)
	start(runGitServer)
}

func initChdir() error {
	root := filepath.Join(os.Getenv("GOPATH"), "src/github.com/growerlab/growerlab/src")
	return os.Chdir(root)
}

func runGitServer() error {
	cfg := configurator.GetConf()
	common.GoSafe(func() {
		err := gggrpc.NewServer(cfg.GitRepoDir, cfg.GoGitGrpcServerAddr)
		if err != nil {
			panic(err)
		}
	})
	time.Sleep(1 * time.Second)
	return nil
}

func start(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}
