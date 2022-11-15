package namespace

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/growerlab/src/backend/app/model/utils"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/jmoiron/sqlx"
)

var table = "namespace"
var columns = []string{
	"id",
	"path",
	"owner_id",
	"type",
}

func AddNamespace(tx sqlx.Queryer, ns *Namespace) error {
	sql, args, _ := sq.Insert(table).
		Columns(columns[1:]...).
		Values(
			ns.Path,
			ns.OwnerID,
			ns.Type,
		).
		Suffix(utils.SqlReturning("id")).
		ToSql()

	err := tx.QueryRowx(sql, args...).Scan(&ns.ID)
	if err != nil {
		return errors.SQLError(err)
	}
	return nil
}

func GetNamespaceByPath(src sqlx.Queryer, path string) (*Namespace, error) {
	return getNamespaceByCond(src, sq.Eq{"path": path})
}

func GetNamespaceByOwnerID(src sqlx.Queryer, ownerID int64) (*Namespace, error) {
	return getNamespaceByCond(src, sq.Eq{"owner_id": ownerID})
}

func GetNamespace(src sqlx.Queryer, id int64) (*Namespace, error) {
	return getNamespaceByCond(src, sq.Eq{"id": id})
}

func getNamespaceByCond(src sqlx.Queryer, cond sq.Sqlizer) (*Namespace, error) {
	ns, err := listNamespaceByCond(src, cond)
	if err != nil {
		return nil, err
	}
	if len(ns) > 0 {
		return ns[0], nil
	}
	return nil, nil
}

func ListNamespacesByOwner(src sqlx.Queryer, userType NamespaceType, ownerIDs ...int64) ([]*Namespace, error) {
	where := sq.And{
		sq.Eq{"owner_id": ownerIDs},
		sq.Eq{"type": userType},
	}
	return listNamespaceByCond(src, where)
}

func ListNamespaceByIDs(src sqlx.Queryer, ids ...int64) ([]*Namespace, error) {
	where := sq.Eq{"id": ids}
	return listNamespaceByCond(src, where)
}

func MapNamespacesByIDs(src sqlx.Queryer, ids ...int64) (map[int64]*Namespace, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	ns, err := ListNamespaceByIDs(src, ids...)
	if err != nil {
		return nil, errors.Trace(err)
	}
	result := make(map[int64]*Namespace)
	for _, n := range ns {
		result[n.ID] = n
	}
	return result, nil
}

func listNamespaceByCond(src sqlx.Queryer, cond sq.Sqlizer) ([]*Namespace, error) {
	sql, args, _ := sq.Select(columns...).From(table).Where(cond).ToSql()

	result := make([]*Namespace, 0)
	err := sqlx.Select(src, &result, sql, args...)
	if err != nil {
		return nil, errors.SQLError(err)
	}
	return result, nil
}
