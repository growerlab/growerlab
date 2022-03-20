package db

import (
	"github.com/growerlab/growerlab/src/backend/app/model/db"
	"github.com/growerlab/growerlab/src/mensa/app/conf"
)

var DB *db.DBQuery

func InitDatabase() error {
	var err error
	var cfg = conf.GetConfig()
	DB, err = db.DoInitDatabase(cfg.DBUrl, conf.IsDev())
	if err != nil {
		return err
	}

	db.DB, err = db.DoInitDatabase(cfg.DBUrl, conf.IsDev())
	return err
}
