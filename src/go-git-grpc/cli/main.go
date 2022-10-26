package main

import (
	"errors"
	"log"
	"os"

	"github.com/growerlab/growerlab/src/common/configurator"
	gggrpc "github.com/growerlab/growerlab/src/go-git-grpc"
)

const (
	hulkPath = "conf/hooks/update"
)

func init() {
	onStart(configurator.InitConfig)
	onStart(checkHook)
}

func main() {
	var err error
	cfg := configurator.GetConf()

	log.Println("go-git-grpc running...", cfg.GoGitGrpcServerAddr)
	log.Println("git root:", cfg.GitRepoDir)

	err = gggrpc.NewServer(cfg.GitRepoDir, cfg.GoGitGrpcServerAddr)
	if err != nil {
		panic(err)
	}
}

func checkHook() error {
	// check git repo hooks
	if _, err := os.Stat(hulkPath); os.IsNotExist(err) {
		return errors.New("git repo hooks not found")
	}
	return nil
}

func onStart(fn func() error) {
	if err := fn(); err != nil {
		log.Printf("%+v\n", err)
		panic(err)
	}
}
