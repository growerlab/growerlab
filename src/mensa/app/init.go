package app

import (
	"io"
	"log"
	"os"

	"github.com/growerlab/growerlab/src/common/configurator"

	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/mensa/app/common"
	"github.com/growerlab/growerlab/src/mensa/app/middleware"
)

// var mids *middleware.Middleware
var manager *Manager
var logger io.Writer = os.Stdout

func initialize() {
	// 初始化日志输出
	log.SetPrefix("[MENSA] ")
	log.SetOutput(logger)

	// 初始化依赖顺序的「初始化」
	startInit(configurator.InitConfig)
	startInit(db.InitDatabase)
	startInit(db.InitMemDB)
	startInit(common.InitPermission)
}

func startInit(fn func() error) {
	if err := fn(); err != nil {
		panic(err)
	}
}

func Run() {
	initialize()

	mensaConf := configurator.GetConf().Mensa

	// 初始化中间件
	entry := new(middleware.Middleware)
	entry.Add(middleware.Authenticate)

	// 初始化管理器
	manager = NewManager(entry)
	manager.RegisterServer(NewGitHttpServer(mensaConf))
	manager.RegisterServer(NewGitSSHServer(mensaConf))
	manager.Run()
}
