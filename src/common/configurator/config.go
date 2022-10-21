package configurator

import (
	"net/url"
	"os"

	"github.com/growerlab/growerlab/src/common/errors"
)

const DefaultConfigPath = "conf/config.yaml"
const DefaultENV = "local"

type Redis struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Namespace   string `yaml:"namespace"`
	MaxIdle     int    `yaml:"max_idle"`
	MaxActive   int    `yaml:"max_active"`
	IdleTimeout int    `yaml:"idle_timeout"`
}

type Mensa struct {
	User        string `yaml:"user"`
	SSHListen   string `yaml:"ssh_listen"`
	HTTPListen  string `yaml:"http_listen"`
	Deadline    int    `yaml:"deadline"`
	IdleTimeout int    `yaml:"idle_timeout"`
	GitPath     string `yaml:"git_path"`
}

type Config struct {
	websiteURL *url.URL

	Debug       bool   `yaml:"debug"`
	WebsiteURL  string `yaml:"website_url"`
	BackendPort int    `yaml:"backend_port"`

	DBUrl string `yaml:"db_url"`
	Redis *Redis `yaml:"redis"`
	Mensa *Mensa `yaml:"mensa"`

	GitRepoDir          string `yaml:"git_repo_dir"`
	GoGitGrpcServerAddr string `yaml:"go_git_grpc_server_addr"`
}

func (c *Config) validate() error {
	if c.Mensa.User == "" {
		return errors.New("git uesr is required")
	}

	if _, err := os.Stat(c.GitRepoDir); os.IsNotExist(err) {
		return errors.Message(err, "git repo dir")
	}
	return nil
}

func (c *Config) EnabledHTTPS() bool {
	if c.websiteURL == nil {
		var err error
		c.websiteURL, err = url.Parse(c.WebsiteURL)
		if err != nil {
			panic(err)
		}
	}
	return c.websiteURL.Scheme == "https"
}

func GetConf() *Config {
	return config
}