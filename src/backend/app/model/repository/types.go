package repository

import (
	"fmt"
	"net"
	"strings"

	"github.com/growerlab/growerlab/src/backend/app/model/namespace"
	"github.com/growerlab/growerlab/src/backend/app/model/user"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/db"
)

type Repository struct {
	ID          int64  `db:"id"`
	UUID        string `db:"uuid"`         // 全站唯一ID（fork时用到）
	Path        string `db:"path"`         // 在namespace中是唯一的name
	Name        string `db:"name"`         // 目前与path字段相同
	NamespaceID int64  `db:"namespace_id"` // 仓库属于个人，还是组织
	OwnerID     int64  `db:"owner_id"`     // 仓库创建者
	Description string `db:"description"`
	CreatedAt   int64  `db:"created_at"`
	ServerID    int64  `db:"server_id"`   // 服务器id
	ServerPath  string `db:"server_path"` // 服务器中的绝对路径
	Public      bool   `db:"public"`      // 公有

	ns    *namespace.Namespace
	owner *user.User
}

// TODO N+1 问题
func (r *Repository) Namespace() *namespace.Namespace {
	if r.ns != nil {
		return r.ns
	}
	r.ns, _ = namespace.GetNamespace(db.DB, r.NamespaceID)
	return r.ns
}

// TODO N+1 问题
func (r *Repository) Owner() *user.User {
	if r.owner != nil {
		return r.owner
	}
	r.owner, _ = user.GetUser(db.DB, r.OwnerID)
	return r.owner
}

func (r *Repository) IsPublic() bool {
	return r.Public
}

func (r *Repository) PathGroup() string {
	return fmt.Sprintf("%s/%s", r.Namespace().Path, r.Path)
}

// https://domain.com:port/user/path.git
func (r *Repository) GitHttpURL() string {
	cfg := configurator.GetConf()

	var sb strings.Builder
	sb.WriteString(cfg.WebsiteURL)
	sb.WriteByte('/')
	sb.WriteString(r.PathGroup())
	sb.WriteString(".git")
	return sb.String()
}

// git@domain.com:port/user/path.git
func (r *Repository) GitSshURL() string {
	cfg := configurator.GetConf().Mensa
	host, rawPort, _ := net.SplitHostPort(cfg.SSHListen)
	port := fmt.Sprintf(":%s", rawPort)
	if rawPort == "22" {
		port = ""
	}

	var sb strings.Builder
	sb.WriteString(cfg.User)
	sb.WriteByte('@')
	sb.WriteString(host)
	sb.WriteString(port)
	sb.WriteByte('/')
	sb.WriteString(r.PathGroup())
	sb.WriteString(".git")
	return sb.String()
}
