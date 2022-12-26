package git

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/test"
	"github.com/growerlab/growerlab/src/go-git-grpc/client"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.InitForTest()
	m.Run()
}

func TestRepository_CreateRepository(t *testing.T) {
	repoPathGroup := "test/admin"
	cfg := configurator.GetConf()
	r := New(context.TODO(), repoPathGroup)
	err := r.Create()
	assert.Nil(t, err)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	// 检测文件
	assert.Equal(t, filepath.Join(cfg.GitRepoDir, "te/ad", repoPathGroup), r.repoAbsPath)
	assert.DirExists(t, r.repoAbsPath)
	assert.FileExists(t, filepath.Join(r.repoAbsPath, "config"))
	assert.DirExists(t, filepath.Join(r.repoAbsPath, "refs"))
}

func TestRepository_DeleteRepository(t *testing.T) {
	repoPathGroup := "test/admin"
	r := New(context.TODO(), repoPathGroup)
	err := r.Delete()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.NoDirExists(t, r.repoAbsPath)
}

func TestRepository_Files(t *testing.T) {
	repoPathGroup := "test/admin"
	r := New(context.TODO(), repoPathGroup)

	file := client.File{
		Ref:         "master",
		AuthorName:  "moli",
		AuthorEmail: "fake@growerlab.net",
		Message:     "test message",
		FilePath:    "test.md",
		FileContent: []byte("# title\n\nhello\n\n"),
	}

	commitHash, err := r.AddFile(&file)
	assert.Nil(t, err)
	assert.NotNil(t, commitHash)
	assert.False(t, plumbing.NewHash(commitHash).IsZero())

	file.FilePath = "test/test2.md"
	file.FileContent = []byte("# tilte\n\ntest2\n\n")
	commitHash, err = r.AddFile(&file)
	assert.Nil(t, err)
	assert.NotNil(t, commitHash)
	assert.False(t, plumbing.NewHash(commitHash).IsZero())

	files, err := r.TreeFiles("master", "/")
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Greater(t, len(files), 1)
}
