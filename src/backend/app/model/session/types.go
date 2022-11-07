package session

import (
	"github.com/growerlab/growerlab/src/backend/app/model/base"
	"github.com/jmoiron/sqlx"
)

const TableName = "session"

var columns = []string{
	"id",
	"owner_id",
	"token",
	"client_ip",
	"created_at",
	"expired_at",
}

type Session struct {
	ID        int64  `db:"id"`
	OwnerID   int64  `db:"owner_id"`
	Token     string `db:"token"`
	ClientIP  string `db:"client_ip"` // 未来可能用来检验token是否被劫持
	CreatedAt int64  `db:"created_at"`
	ExpiredAt int64  `db:"expired_at"`
}

type model struct {
	*base.Model
	src sqlx.Ext
}

func New(src sqlx.Ext) *model {
	return &model{
		src:   src,
		Model: base.NewModel(src, TableName, nil),
	}
}
