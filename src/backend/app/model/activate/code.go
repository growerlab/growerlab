package activate

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/growerlab/src/backend/app/common/errors"
	"github.com/jmoiron/sqlx"
)

var tableName = "activate_code"
var columns = []string{
	"id",
	"user_id",
	"code",
	"created_at",
	"used_at",
	"expired_at",
}

func AddCode(tx sqlx.Execer, code *ActivationCode) error {
	code.CreatedAt = time.Now().Unix()

	sql, args, _ := sq.Insert(tableName).
		Columns(columns[1:]...).
		Values(
			code.UserID,
			code.Code,
			code.CreatedAt,
			nil,
			code.ExpiredAt,
		).ToSql()

	_, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.SQLError(err)
	}
	return nil
}

func GetCode(src sqlx.Queryer, code string) (*ActivationCode, error) {
	sql, args, _ := sq.Select(columns...).
		From(tableName).
		Where(sq.Eq{"code": code}).
		Limit(1).
		ToSql()

	var data = make([]*ActivationCode, 0)
	err := sqlx.Select(src, &data, sql, args...)
	if err != nil {
		return nil, errors.SQLError(err)
	}
	if len(data) > 0 {
		return data[0], nil
	}
	return nil, nil
}

// ActivateCode
func ActivateCode(tx sqlx.Execer, code string) error {
	sql, args, _ := sq.Update(tableName).
		Set("used_at", time.Now().Unix()).
		Where(sq.Eq{"code": code}).
		ToSql()

	_, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.SQLError(err)
	}
	return nil
}
