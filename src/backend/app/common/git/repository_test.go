package git

import (
	"context"
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
	cfg := configurator.GetConf()
	r := &Repository{
		ctx:       context.Background(),
		cfg:       cfg,
		pathGroup: "test/admin",
	}
	err := r.CreateRepository()
	assert.Nil(t, err)

	// TODO 检测文件
}
