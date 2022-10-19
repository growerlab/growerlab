package main

import (
	"log"
	"os"
	"path/filepath"

	gggrpc "github.com/growerlab/growerlab/src/go-git-grpc"
	"gopkg.in/yaml.v3"
)

const (
	defaultConfig = "conf/config.yaml"
	hulkPath      = "conf/hooks/update"
)

type Config struct {
	Listen     string `yaml:"listen"`
	GitRepoDir string `yaml:"git_repo_dir"`
}

var conf *Config

func init() {
	var envConfig = map[string]*Config{}
	var ok bool
	rawConfig, err := os.ReadFile(defaultConfig)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(rawConfig, &envConfig)
	if err != nil {
		panic(err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	if conf, ok = envConfig[env]; !ok {
		panic("env " + env + " not found")
	}

	// check git repo hooks
	if _, err = os.Stat(hulkPath); os.IsNotExist(err) {
		panic("git repo hooks not found")
	}
}

func main() {
	var err error
	root := conf.GitRepoDir
	root, err = filepath.Abs(root)
	if err != nil {
		panic(err)
	}

	log.Println("go-git-grpc running...", conf.Listen)
	log.Println("git root:", root)

	err = gggrpc.NewServer(root, conf.Listen)
	if err != nil {
		panic(err)
	}
}
