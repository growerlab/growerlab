package common

import (
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/permission"
)

func InitPermission() error {
	err := permission.InitPermissionHub(db.DB, db.MemDB)
	return err
}
