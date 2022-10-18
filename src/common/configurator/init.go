package configurator

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/growerlab/growerlab/src/common/errors"
	"gopkg.in/yaml.v3"
)

var (
	config *Config
)

func InitConfig() error {
	var foundEnv bool
	var envConfig = make(map[string]*Config)

	confBody, err := ioutil.ReadFile(DefaultConfigPath)
	if err != nil {
		return errors.Trace(err)
	}
	err = yaml.Unmarshal(confBody, &envConfig)
	if err != nil {
		return errors.Trace(err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = DefaultENV
	}

	config, foundEnv = envConfig[env]
	if !foundEnv {
		return errors.Errorf("config for env '%s'", env)
	}

	err = config.validate()
	if err != nil {
		return errors.Trace(err)
	}
	log.Println("Loaded config for env:", env)
	return err
}
