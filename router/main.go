/*
	dev env.
	nginx -> router -> growerlab

	接受nginx的请求，将请求路由到不同的后端 growerlab
*/

package main

import (
	"github.com/growerlab/growerlab/router/app"
)

func main() {
	app.RunWebSiteRouter()

	app.RunSSHRouter()
	app.RunGitHTTPRouter()
}
