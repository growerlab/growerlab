package common

import (
	"github.com/growerlab/growerlab/src/backend/app/common/permission"
	"github.com/growerlab/growerlab/src/common/db"
)

func InitPermission() error {
	err := permission.InitPermissionHub(db.DB, db.MemDB)
	return err
}
