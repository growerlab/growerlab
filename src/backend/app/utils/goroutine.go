package utils

import (
	"log"
)

func GoSafe(f func()) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("go safe: %+v", e)
		}
	}()
	go f()
}
