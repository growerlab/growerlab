package app

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// 读取宿主的配置文件（例如go-git-grpc）
const configPath = "conf/config.yaml"

var Conf *Config

type RedisConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	MaxIdle     int    `yaml:"max_idle"`
	MaxActive   int    `yaml:"max_active"`
	IdleTimeout int    `yaml:"idle_timeout"`
	Namespace   string `yaml:"namespace"`
}

type Config struct {
	Redis *RedisConfig `yaml:"redis"`
}

func InitConfig() error {
	var ok bool
	var env = os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	confRaw, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println(os.Getwd())
		panic(err)
	}

	envSet := make(map[string]*Config)

	err = yaml.Unmarshal(confRaw, &envSet)
	if err != nil {
		return errors.WithStack(err)
	}
	if Conf, ok = envSet[env]; !ok {
		panic("config not found")
	}
	return nil
}
