package git

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/test"
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

	assert.Nil(t, err)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.NoDirExists(t, r.repoAbsPath)
}
