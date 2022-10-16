package main

import (
	"fmt"
	"log"
	"os"

	"github.com/growerlab/growerlab/src/mensa/script/hulk/app"
)

func main() {
	l := app.NewLogger(fmt.Sprintf("%s/%s", app.EnvRepoOwner, app.EnvRepoPath))
	log.SetOutput(l)

	defer func() {
		l.Flush()
		if e := recover(); e != nil {
			log.Println(e)
			os.Exit(1)
		}
	}()

	sess := app.Session()
	if err := app.Run(sess); err != nil {
		app.ErrPanic(err)
	}
}
