package user

import (
	"github.com/growerlab/growerlab/src/backend/app/model/namespace"
	"github.com/growerlab/growerlab/src/common/db"
)

var tableName = "user"
var tableNameMark = "`user`" // user 是 pgsql中的保留关键字，所以加上引号

var columns = []string{
	"id",
	"email",
	"encrypted_password",
	"username",
	"name",
	"public_email",
	"created_at",
	"deleted_at",
	"verified_at",
	"last_login_at",
	"last_login_ip",
	"register_ip",
	"is_admin",
	"namespace_id",
}

type User struct {
	ID                int64   `db:"id"`
	Email             string  `db:"email"`
	EncryptedPassword string  `db:"encrypted_password"`
	Username          string  `db:"username"`
	Name              string  `db:"name"`
	PublicEmail       string  `db:"public_email"`
	CreatedAt         int64   `db:"created_at"`
	DeletedAt         *int64  `db:"deleted_at"`
	VerifiedAt        *int64  `db:"verified_at"`
	LastLoginAt       *int64  `db:"last_login_at"`
	LastLoginIP       *string `db:"last_login_ip"`
	RegisterIP        string  `db:"register_ip"`
	IsAdmin           bool    `db:"is_admin"`
	NamespaceID       int64   `db:"namespace_id"`

	ns *namespace.Namespace // cached namespace
}

// TODO N+1 问题
func (u *User) Namespace() *namespace.Namespace {
	if u.ns != nil {
		return u.ns
	}
	u.ns, _ = namespace.GetNamespaceByOwnerID(db.DB, u.ID)
	return u.ns
}

func (u *User) Verified() bool {
	return u.VerifiedAt != nil && *u.VerifiedAt > 0
}
