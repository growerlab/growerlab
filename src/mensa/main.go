package main

import (
	"fmt"

	"github.com/growerlab/growerlab/src/mensa/app"
)

func main() {
	fmt.Println("=================================")
	fmt.Println("BuiltTime: ", app.BUILDTIME)
	fmt.Println("CommitID: ", app.BUILDCOMMIT)
	fmt.Println("GoVersion: ", app.GOVERSION)
	fmt.Println("=================================")
	fmt.Println(app.UA)

	app.Run()
}
