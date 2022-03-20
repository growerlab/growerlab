package repo

import (
	"github.com/go-git/go-git/v5/plumbing"
)

func IsBranch(ref string) bool {
	return plumbing.ReferenceName(ref).IsBranch()
}

func IsTag(ref string) bool {
	return plumbing.ReferenceName(ref).IsTag()
}
