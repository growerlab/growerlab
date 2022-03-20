package session

import (
	"github.com/growerlab/growerlab/src/backend/app/common/errors"
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

func (m *model) Add(sess *Session) error {
	values := []interface{}{
		sess.OwnerID,
		sess.Token,
		sess.ClientIP,
		sess.CreatedAt,
		sess.ExpiredAt,
	}
	var err error
	sess.ID, err = m.Insert(columns[1:], values).Exec()
	return errors.SQLError(err)
}
