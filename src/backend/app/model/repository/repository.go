package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/growerlab/src/backend/app/common/errors"
	"github.com/growerlab/growerlab/src/backend/app/model/utils"
	"github.com/jmoiron/sqlx"
)

var (
	table   = "repository"
	columns = []string{
		"id",
		"uuid",
		"path",
		"name",
		"namespace_id",
		"owner_id",
		"description",
		"created_at",
		"server_id",
		"server_path",
		"public",
	}
)

func AddRepository(tx sqlx.Ext, repo *Repository) error {
	sql, args, _ := sq.Insert(table).
		Columns(columns[1:]...).
		Values(
			repo.UUID,
			repo.Path,
			repo.Name,
			repo.NamespaceID,
			repo.OwnerID,
			repo.Description,
			repo.CreatedAt,
			repo.ServerID,
			repo.ServerPath,
			repo.Public,
		).
		Suffix(utils.SqlReturning("id")).
		ToSql()

	err := tx.QueryRowx(sql, args...).Scan(&repo.ID)
	if err != nil {
		return errors.SQLError(err)
	}
	return nil
}

func NameExistInNamespace(src sqlx.Queryer, namespaceID int64, name string) (bool, error) {
	where := sq.And{
		sq.Eq{"namespace_id": namespaceID},
		sq.Eq{"path": name},
	}
	result, err := listRepositoriesByCond(src, []string{"id"}, where)
	if err != nil {
		return false, err
	}
	return len(result) > 0, nil
}

func ListRepositoriesByNamespace(src sqlx.Queryer, namespaceID int64) ([]*Repository, error) {
	where := sq.And{sq.Eq{"namespace_id": namespaceID}}
	return listRepositoriesByCond(src, columns, where)
}

func GetRepositoryByNsWithPath(src sqlx.Queryer, namespaceID int64, path string) (*Repository, error) {
	where := sq.And{sq.Eq{"namespace_id": namespaceID, "path": path}}
	repos, err := listRepositoriesByCond(src, columns, where)
	if err != nil {
		return nil, err
	}
	if len(repos) > 0 {
		return repos[0], nil
	}
	return nil, nil
}

func GetRepository(src sqlx.Queryer, id int64) (*Repository, error) {
	repos, err := listRepositoriesByCond(src, columns, sq.Eq{"id": id})
	if err != nil {
		return nil, err
	}
	if len(repos) > 0 {
		return repos[0], nil
	}
	return nil, nil
}

func listRepositoriesByCond(src sqlx.Queryer, tableColumns []string, cond sq.Sqlizer) ([]*Repository, error) {
	where := cond
	sql, args, _ := sq.Select(tableColumns...).
		From(table).
		Where(where).
		ToSql()

	result := make([]*Repository, 0)
	err := sqlx.Select(src, &result, sql, args...)
	if err != nil {
		return nil, errors.SQLError(err)
	}
	return result, nil
}
