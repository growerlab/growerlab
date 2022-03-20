package permission

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/growerlab/src/backend/app/common/context"
	"github.com/growerlab/growerlab/src/backend/app/common/errors"
	"github.com/jmoiron/sqlx"
)

var table = "permission"
var columns = []string{
	"id",
	"namespace_id",
	"code",
	"context_type",
	"context_param_1",
	"context_param_2",
	"user_domain_type",
	"user_domain_param",
	"created_at",
	"deleted_at",
}

func ListPermissionsByContext(src sqlx.Queryer, code int, c *context.Context) ([]*Permission, error) {
	where := sq.And{
		sq.Eq{"code": code},
		sq.Eq{"context_type": c.Type},
		sq.Eq{"context_param_1": c.Param1},
		sq.Eq{"context_param_2": c.Param2},
	}
	return listPermissionByCond(src, columns, where)
}

func listPermissionByCond(src sqlx.Queryer, cols []string, cond sq.Sqlizer) ([]*Permission, error) {
	sql, args, _ := sq.Select(cols...).
		From(table).
		Where(cond).
		ToSql()

	result := make([]*Permission, 0)
	err := sqlx.Select(src, &result, sql, args...)
	if err != nil {
		return nil, errors.SQLError(err)
	}
	return result, nil
}
