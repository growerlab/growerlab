package permission

import (
	"testing"

	"github.com/growerlab/growerlab/src/backend/test"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/stretchr/testify/assert"
)

func init() {
	test.InitDIR()

	onStart(configurator.InitConfig)
	onStart(db.InitMemDB)
	onStart(db.InitDatabase)
	onStart(InitPermission)
}

func onStart(fn func() error) {
	if err := fn(); err != nil {
		panic(err)
	}
}

func TestCheckViewRepository(t *testing.T) {
	usr, repo, err := test.MakeTestRepoData()
	if !assert.Equal(t, nil, err, nil) {
		return
	}
	if !assert.NotEqual(t, nil, usr, nil) {
		return
	}
	if !assert.NotEqual(t, nil, repo, nil) {
		return
	}

	err = CheckViewRepository(&usr.NamespaceID, repo.ID)
	assert.Equal(t, nil, err, nil) // has permission

	err = CheckCloneRepository(&usr.NamespaceID, repo.ID)
	assert.Equal(t, nil, err, nil) // has permission

	err = CheckPushRepository(usr.NamespaceID, repo.ID)
	assert.Equal(t, nil, err, nil) // has permission

	err = CheckPushRepository(0, repo.ID) // err user
	assert.NotEqual(t, nil, err, nil)     // no permission
}
