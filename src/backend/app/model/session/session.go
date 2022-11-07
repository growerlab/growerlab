package session

import (
	"github.com/growerlab/growerlab/src/common/errors"
)

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
