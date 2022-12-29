package main

import (
	"fmt"
	"net/http"

	"github.com/growerlab/growerlab/src/backend/app"
	"github.com/growerlab/growerlab/src/backend/app/common/notify"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/logger"
)

func main() {
	addr := fmt.Sprintf(":%d", configurator.GetConf().BackendPort)
	err := app.Run(addr)
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info("bye.")
		} else {
			panic(err)
		}
	}
	notify.WaitDone()
}
