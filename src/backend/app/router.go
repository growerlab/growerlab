package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/backend/app/common/notify"
	"github.com/growerlab/growerlab/src/backend/app/controller"
	"github.com/growerlab/growerlab/src/backend/app/controller/middleware"
	"github.com/growerlab/growerlab/src/common/configurator"
)

func Run(addr string) error {
	engine := gin.Default()

	engine.Use(middleware.CORSForLocal)
	apiV1 := engine.Group("/api/v1", middleware.LimitGETRequestBody)

	repositories := apiV1.Group("/repositories")
	{
		repositories.POST("/:namespace/create", middleware.Authorized, controller.CreateRepository)
		repositories.GET("/:namespace/list", controller.Repositories)
		repositories.GET("/:namespace/detail/:name", controller.Repository)
	}

	auth := apiV1.Group("/auth")
	{
		auth.POST("/register", controller.RegisterUser)
		auth.POST("/activate", controller.ActivateUser)
		auth.POST("/login", controller.LoginUser)
	}

	return runServer(addr, engine)
}

func runServer(addr string, engine *gin.Engine) error {
	server := &http.Server{
		Addr:         addr,
		Handler:      engine,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	// 平滑关闭
	notify.Subscribe(func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		server.Shutdown(timeoutCtx)
	})

	// 是否debug
	gin.SetMode(gin.ReleaseMode)
	if configurator.GetConf().Debug {
		gin.SetMode(gin.DebugMode)
	}
	return server.ListenAndServe()
}
