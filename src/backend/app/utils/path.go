package utils

import (
	"path/filepath"
)

func BasePath() string {
	dir, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}
	return dir
}
