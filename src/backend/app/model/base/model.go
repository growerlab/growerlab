package base

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/growerlab/src/backend/app/model/utils"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/jmoiron/sqlx"
)

type Model struct {
	src sqlx.Ext

	table string
	alias string

	defaultTerms       sq.And
	ignoreDefaultTerms bool // 是否忽略默认条件
}

func NewModel(src sqlx.Ext, table string, defaultTerms sq.And) *Model {
	return &Model{
		src:          src,
		table:        table,
		defaultTerms: defaultTerms,
	}
}

func (m *Model) Alias(a string) *Model {
	m.alias = a
	return m
}

func (m *Model) IgnoreDefaultTerms() {
	m.ignoreDefaultTerms = true
}

func (m *Model) Select() *Selector {
	where := sq.And{}
	if m.defaultTerms != nil && !m.ignoreDefaultTerms {
		where = append(where, m.defaultTerms...)
	}

	table := m.getTable()
	builder := sq.Select().From(table).Where(where)
	return &Selector{
		SelectBuilder: builder,
		src:           m.src,
		tableName:     table,
	}
}

func (m *Model) getTable() string {
	if len(m.alias) > 0 {
		return fmt.Sprintf("%s AS %s", m.table, m.alias)
	}
	return m.table
}

func (m *Model) Update(set map[string]interface{}, term sq.Sqlizer) *Updater {
	if len(set) == 0 {
		panic("'set' must required")
	}

	where := sq.And{}
	where = append(where, term)

	builder := sq.Update(m.table).SetMap(set).Where(where)

	return &Updater{
		src:     m.src,
		table:   m.table,
		values:  set,
		builder: builder,
	}
}

func (m *Model) Delete(term sq.Sqlizer) *Deleter {
	where := sq.And{}
	where = append(where, term)

	// 提取出删除条件
	values := make(map[string]interface{})
	switch term.(type) {
	case sq.And:
		ands := term.(sq.And)
		for i := range ands {
			if eq, ok := ands[i].(sq.Eq); ok {
				for k, v := range eq {
					values[k] = v
				}
			}
		}
	case sq.Eq:
		eq := term.(sq.Eq)
		for k, v := range eq {
			values[k] = v
		}
	}

	builder := sq.Delete(m.table).Where(where)

	return &Deleter{
		src:     m.src,
		table:   m.table,
		values:  values,
		builder: builder,
	}
}

type Deleter struct {
	src     sqlx.Ext
	table   string
	values  map[string]interface{}
	builder sq.DeleteBuilder
}

func (d *Deleter) Exec() error {
	// hook before
	err := hook.Effect(d.src, d.table, ActionDelete, TenseBefore, d.values)
	if err != nil {
		return err
	}

	_, err = d.builder.RunWith(d.src).Exec()
	if err != nil {
		return errors.SQLError(err)
	}

	// hook after
	err = hook.Effect(d.src, d.table, ActionDelete, TenseAfter, d.values)
	return err
}

func (m *Model) Insert(columns []string, values []interface{}) *Inserter {
	if len(values) == 0 || len(columns) == 0 {
		panic("'columns' and 'values' must required")
	}

	builder := sq.Insert(m.table).Columns(columns...)

	return &Inserter{
		src:     m.src,
		table:   m.table,
		columns: columns,
		values:  values,
		builder: builder,
	}
}

// BatchInsert 批量插入，不会触发hook
func (m *Model) BatchInsert(columns []string, size int, getValuesFn func(int) []interface{}) error {
	const maxValues = 1000
	valueBucket := make([][]interface{}, 0, maxValues)

	batchInsertFunc := func(mulValues [][]interface{}) error {
		builder := sq.Insert(m.table).Columns(columns...)
		for _, values := range mulValues {
			builder = builder.Values(values...)
		}
		_, err := builder.RunWith(m.src).Exec()
		return errors.SQLError(err)
	}

	for i := 0; i < size; i++ {
		values := getValuesFn(i)
		valueBucket = append(valueBucket, values)
		if len(valueBucket) >= maxValues || i == size-1 {
			if err := batchInsertFunc(valueBucket); err != nil {
				return err
			}
			valueBucket = valueBucket[:0]
		}
	}

	return nil
}

type Inserter struct {
	src     sqlx.Ext
	table   string
	columns []string
	values  []interface{}
	builder sq.InsertBuilder

	sqlReturningColumn string
}

func (i *Inserter) SqlReturning(column string) db.Execor {
	i.builder = i.builder.Suffix(utils.SqlReturning(column))
	i.sqlReturningColumn = column
	return i
}

func (i *Inserter) Exec() (int64, error) {
	// 单个插入数据
	set := make(map[string]interface{})
	for v := range i.values {
		set[i.columns[v]] = i.values[v]
	}
	// hook before
	err := hook.Effect(i.src, i.table, ActionCreate, TenseBefore, set)
	if err != nil {
		return -1, err
	}

	i.builder = i.builder.Values(i.values...)
	query, args, err := i.builder.ToSql()
	if err != nil {
		return -1, errors.Trace(err)
	}
	row := i.src.QueryRowx(query, args...)

	// hook after
	err = hook.Effect(i.src, i.table, ActionCreate, TenseAfter, set)
	if err != nil {
		return -1, errors.Trace(err)
	}

	if len(i.sqlReturningColumn) > 0 {
		var returningVal int64
		err = row.Scan(&returningVal)
		if err != nil {
			return -1, errors.Trace(err)
		}
		return returningVal, nil
	}

	return -1, nil
}

type Updater struct {
	src     sqlx.Ext
	table   string
	values  map[string]interface{}
	builder sq.UpdateBuilder
}

func (u *Updater) Exec() error {
	// hook before
	err := hook.Effect(u.src, u.table, ActionUpdate, TenseBefore, u.values)
	if err != nil {
		return errors.Trace(err)
	}

	_, err = u.builder.RunWith(u.src).Exec()
	if err != nil {
		return errors.SQLError(err)
	}

	// hook after
	err = hook.Effect(u.src, u.table, ActionUpdate, TenseAfter, u.values)
	return errors.Trace(err)
}

type Selector struct {
	sq.SelectBuilder
	src       sqlx.Ext
	tableName string
}

func (s *Selector) BuildSQL(fn func(builder sq.SelectBuilder) sq.SelectBuilder) *Selector {
	// s.builder = fn(s.builder)
	panic("not support")
	// return s
}

func (s *Selector) Query(dest interface{}) error {
	query, args, err := s.ToSql()
	if err != nil {
		return errors.Trace(err)
	}
	err = sqlx.Select(s.src, dest, query, args...)
	return errors.SQLError(err)
}
