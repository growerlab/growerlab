package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/jmoiron/sqlx"
)

func (m *model) AddRepository(repo *Repository) error {
	values := []interface{}{
		repo.UUID,
		repo.Path,
		repo.Name,
		repo.NamespaceID,
		repo.OwnerID,
		repo.Description,
		repo.CreatedAt,
		repo.Public,
	}
	_, err := m.Insert(columns[1:], values).Exec()
	return errors.SQLError(err)
}

func (m *model) NameExists(namespaceID int64, name string) (bool, error) {
	where := sq.And{
		sq.Eq{"namespace_id": namespaceID},
		sq.Eq{"path": name},
	}
	result, err := m.listRepositoriesByCond([]string{"id"}, where)
	if err != nil {
		return false, err
	}
	return len(result) > 0, nil
}

func (m *model) ListRepositoriesByNamespace(namespaceID int64) ([]*Repository, error) {
	where := sq.And{sq.Eq{"namespace_id": namespaceID}}
	return m.listRepositoriesByCond(columns, where)
}

func (m *model) GetRepositoryByNsWithPath(namespaceID int64, path string) (*Repository, error) {
	where := sq.And{sq.Eq{"namespace_id": namespaceID, "path": path}}
	repos, err := m.listRepositoriesByCond(columns, where)
	if err != nil {
		return nil, err
	}
	if len(repos) > 0 {
		return repos[0], nil
	}
	return nil, nil
}

func (m *model) GetRepository(id int64) (*Repository, error) {
	repos, err := m.listRepositoriesByCond(columns, sq.Eq{"id": id})
	if err != nil {
		return nil, err
	}
	if len(repos) > 0 {
		return repos[0], nil
	}
	return nil, nil
}

func (m *model) listRepositoriesByCond(tableColumns []string, cond sq.Sqlizer) ([]*Repository, error) {
	where := cond
	sql, args, _ := sq.Select(tableColumns...).
		From(tableName).
		Where(where).
		ToSql()

	result := make([]*Repository, 0)
	err := sqlx.Select(m.src, &result, sql, args...)
	if err != nil {
		return nil, errors.SQLError(err)
	}
	return result, nil
}
