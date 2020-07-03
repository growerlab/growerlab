package app

import (
	"errors"
	"strings"
)

func BranchByHost(host string) string {
	if !strings.Contains(host, ".dev.growerlab.net") {
		panic(errors.New("invalid host"))
	}
	n := strings.Index(host, ".")
	return host[:n]
}
