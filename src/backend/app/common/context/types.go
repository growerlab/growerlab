package context

import (
	"github.com/growerlab/growerlab/src/backend/app/model/db"
	"github.com/jmoiron/sqlx"
)

type Context struct {
	Type   int   `json:"type"`
	Param1 int64 `json:"param1"`
	Param2 int64 `json:"param2"`
}

type DBContext struct {
	Src   sqlx.Queryer
	MemDB *db.MemDBClient
}
