package user

import (
	"fmt"
	"github.com/growerlab/growerlab/src/common/errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/growerlab/src/backend/app/model/namespace"
	"github.com/growerlab/growerlab/src/backend/app/model/session"
	"github.com/growerlab/growerlab/src/backend/app/model/utils"
	"github.com/jmoiron/sqlx"
)

func AddUser(tx sqlx.Queryer, user *User) error {
	sql, args, _ := sq.Insert(tableNameMark).
		Columns(columns[1:]...).
		Values(
			user.Email,
			user.EncryptedPassword,
			user.Username,
			user.Name,
			user.PublicEmail,
			user.CreatedAt,
			nil,
			nil,
			nil,
			nil,
			user.RegisterIP,
			user.IsAdmin,
			user.NamespaceID,
		).
		Suffix(utils.SqlReturning("id")).
		ToSql()

	err := tx.QueryRowx(sql, args...).Scan(&user.ID)
	if err != nil {
		return errors.SQLError(err)
	}
	return nil
}

func ExistsEmailOrUsername(src sqlx.Queryer, username, email string) (bool, error) {
	if len(username) > 0 {
		user, err := getUser(src, sq.Eq{"username": username})
		if err != nil {
			return false, err
		}
		if user != nil {
			return true, nil
		}
	}
	if len(email) > 0 {
		user, err := getUser(src, sq.Eq{"email": email})
		if err != nil {
			return false, err
		}
		if user != nil {
			return true, nil
		}
	}
	return false, nil
}

func GetUserByEmail(src sqlx.Queryer, email string) (*User, error) {
	user, err := getUser(src, sq.Eq{"email": email})
	return user, err
}

func GetUserByUsername(src sqlx.Queryer, username string) (*User, error) {
	user, err := getUser(src, sq.Eq{"username": username})
	return user, err
}

func GetUser(src sqlx.Queryer, id int64) (*User, error) {
	user, err := getUser(src, sq.Eq{"id": id})
	return user, err
}

func getUser(src sqlx.Queryer, cond sq.Sqlizer) (*User, error) {
	users, err := listUsersByCond(src, columns, cond)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return users[0], nil
	}
	return nil, nil
}

func listUsersByCond(src sqlx.Queryer, tableColumns []string, cond sq.Sqlizer) ([]*User, error) {
	sql, args, _ := sq.Select(tableColumns...).
		From(tableNameMark).
		Where(sq.And{cond, NormalUser}).
		ToSql()

	result := make([]*User, 0)
	err := sqlx.Select(src, &result, sql, args...)
	if err != nil {
		return nil, errors.SQLError(err)
	}
	return result, nil
}

func ActivateUser(tx sqlx.Execer, userID int64) error {
	sql, args, _ := sq.Update(tableNameMark).
		Set("verified_at", time.Now().Unix()).
		Where(sq.And{sq.Eq{"id": userID}, InactivateUser}).
		ToSql()

	_, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.SQLError(err)
	}
	return nil
}

func ListAllUsers(src sqlx.Queryer, page, per uint64) ([]*User, error) {
	users := make([]*User, 0)

	// TODO 如果用户量很大的时候，这样分页会有性能问题.. 希望能碰到那一天👀
	sql, _, _ := sq.Select(columns...).
		From(tableNameMark).
		Where(NormalUser).
		Limit(per).
		Offset(page * per).
		ToSql()

	err := sqlx.Select(src, &users, sql)
	return users, errors.SQLError(err)
}

func UpdateLogin(tx sqlx.Execer, userID int64, clientIP string) error {
	where := sq.Eq{"id": userID}
	valueMap := map[string]interface{}{
		"last_login_at": time.Now().Unix(),
		"last_login_ip": clientIP,
	}
	return update(tx, where, valueMap)
}

func UpdateNamespace(tx sqlx.Execer, userID int64, namespaceID int64) error {
	where := sq.Eq{"id": userID}
	valueMap := map[string]interface{}{
		"namespace_id": namespaceID,
	}
	return update(tx, where, valueMap)
}

func update(tx sqlx.Execer, cond sq.Sqlizer, valueMap map[string]interface{}) error {
	sql, args, _ := sq.Update(tableNameMark).
		SetMap(valueMap).
		Where(cond).
		ToSql()

	_, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.SQLError(err)
	}
	return nil
}

func GetUserByUserToken(src sqlx.Queryer, userToken string) (*User, error) {
	sessTableName := session.TableName
	joinColumns := utils.SqlColumnsComplementTable(tableNameMark, columns...)
	sql, args, _ := sq.Select(joinColumns...).
		From(tableNameMark).
		Join(fmt.Sprintf("%s ON %s.token = ? AND %s.expired_at >= ?", sessTableName, sessTableName, sessTableName),
			userToken, time.Now().Unix()).
		Where(fmt.Sprintf("%s.id = %s.owner_id", tableNameMark, sessTableName)).
		ToSql()

	users := make([]*User, 0, 1)

	err := sqlx.Select(src, &users, sql, args...)
	if err != nil {
		return nil, errors.SQLError(err)
	}
	if len(users) > 0 {
		return users[0], nil
	}
	return nil, nil
}

func ListAdminUsers(src sqlx.Queryer) ([]*User, error) {
	where := sq.And{
		sq.Eq{"is_admin": true},
	}
	users, err := listUsersByCond(src, columns, where)
	if err != nil {
		return nil, err
	}
	err = fillNamespaceInUsers(src, users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func fillNamespaceInUsers(src sqlx.Queryer, users []*User) error {
	userIDs := make([]int64, 0)
	userMap := make(map[int64]*User)
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
		userMap[u.ID] = u
	}
	ns, err := namespace.ListNamespacesByOwner(src, namespace.TypeUser, userIDs...)
	if err != nil {
		return err
	}
	// fill
	for _, n := range ns {
		userMap[n.OwnerID].ns = n
	}
	return nil
}
