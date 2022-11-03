package git

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/path"
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
	r := &Repository{
		ctx:       context.Background(),
		cfg:       cfg,
		pathGroup: repoPathGroup,
	}
	err := r.CreateRepository()
	assert.Nil(t, err)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	// 检测文件
	absolutePath := path.GetRealRepositoryPath(repoPathGroup)
	assert.Equal(t, filepath.Join(cfg.GitRepoDir, "te/ad", repoPathGroup), absolutePath)
	assert.DirExists(t, absolutePath)
	assert.FileExists(t, filepath.Join(absolutePath, "config"))
	assert.DirExists(t, filepath.Join(absolutePath, "refs"))
}
