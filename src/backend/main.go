package main

import (
	"fmt"
	"net/http"

	"github.com/growerlab/growerlab/src/backend/app"
	"github.com/growerlab/growerlab/src/backend/app/common/notify"
	"github.com/growerlab/growerlab/src/backend/app/utils/conf"
	"github.com/growerlab/growerlab/src/backend/app/utils/logger"
)

func main() {
	addr := fmt.Sprintf(":%d", conf.GetConf().Port)
	err := app.Run(addr)
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info("bye.")
		} else {
			panic(err)
		}
	}
	notify.Done()
}
