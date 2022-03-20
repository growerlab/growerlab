// KeyDB / Redis 配置

package db

import (
	"github.com/growerlab/growerlab/src/backend/app/model/db"
	selfConf "github.com/growerlab/growerlab/src/mensa/app/conf"
)

var MemDB *db.MemDBClient
var PermissionDB *db.MemDBClient

func InitMemDB() (err error) {
	cfg := selfConf.GetConfig()
	redisConf := cfg.Redis
	defaultRedisConf := redisConf.Redis

	permissionRedisConf := redisConf.Redis
	permissionRedisConf.Namespace = redisConf.PermissionNamespace

	MemDB, err = db.DoInitMemDB(&defaultRedisConf, 0)
	if err != nil {
		return err
	}

	PermissionDB, err = db.DoInitMemDB(&permissionRedisConf, 0)
	if err != nil {
		return err
	}
	return nil
}
