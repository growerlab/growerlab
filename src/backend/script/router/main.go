/*
	dev env.
	nginx -> router -> growerlab

	接受nginx的请求，将请求路由到不同的后端 growerlab
*/

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/growerlab/growerlab/src/backend/script/router/app"
)

func main() {
	app.RunWebSiteRouter()

	app.RunSSHRouter()
	app.RunGitHTTPRouter()

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	<-c
}
