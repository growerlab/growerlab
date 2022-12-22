package command

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

var root = filepath.Join(os.Getenv("GOPATH"), "src/github.com/growerlab/growerlab/.repositories")

func TestBlob_Commit(t *testing.T) {
	b := NewBlob(
		root,
		"dir/test.md", []byte("## title"), &Context{
			Bin:      "git",
			RepoPath: "te/ad/test/admin",
		})

	got, err := b.Commit(object.Signature{
		Name:  "moli",
		Email: "moli@admin.com",
		When:  time.Now(),
	}, "test\nhello", "main")

	assert.Nil(t, err, nil)
	assert.NotEqual(t, got, plumbing.ZeroHash)
}
