package common

import (
	"github.com/growerlab/growerlab/src/backend/app/common/permission"
	"github.com/growerlab/growerlab/src/mensa/app/db"
)

func InitPermission() error {
	err := permission.InitPermissionHub(db.DB, db.PermissionDB)
	return err
}
