package test

import (
	"github.com/growerlab/growerlab/src/common/configurator"
)

func InitForTest() {
	start(configurator.InitConfig)
}

func start(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}
