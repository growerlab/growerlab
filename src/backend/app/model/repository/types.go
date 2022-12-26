package repository

import (
	"fmt"
	"net"
	"strings"

	"github.com/growerlab/growerlab/src/backend/app/model/base"
	"github.com/growerlab/growerlab/src/backend/app/model/namespace"
	"github.com/growerlab/growerlab/src/backend/app/model/user"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/path"
	"github.com/jmoiron/sqlx"
)

var (
	tableName = "repository"
	columns   = []string{
		"id",
		"uuid",
		"path",
		"name",
		"namespace_id",
		"owner_id",
		"description",
		"created_at",
		"public",
		"last_push_at",
		"default_branch",
	}
)

type Repository struct {
	ID            int64  `db:"id"`
	UUID          string `db:"uuid"`         // 全站唯一ID（fork时用到）
	Path          string `db:"path"`         // 在namespace中是唯一的name
	Name          string `db:"name"`         // 目前与path字段相同
	NamespaceID   int64  `db:"namespace_id"` // 仓库所属的命名空间（个人，组织）
	OwnerID       int64  `db:"owner_id"`     // 仓库创建者
	Description   string `db:"description"`
	CreatedAt     int64  `db:"created_at"`
	Public        bool   `db:"public"`         // 公有
	LastPushAt    int64  `db:"last_push_at"`   // 最后的推送时间
	DefaultBranch string `db:"default_branch"` // 默认分支

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
	return path.GetPathGroup(r.Namespace().Path, r.Path)
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

func FillNamespaces(tx sqlx.Ext, repos ...*Repository) error {
	if len(repos) == 0 {
		return nil
	}

	nsIDs := make([]int64, 0, len(repos))
	for _, repo := range repos {
		nsIDs = append(nsIDs, repo.NamespaceID)
	}
	namespaceSet, err := namespace.MapNamespacesByIDs(tx, nsIDs...)
	if err != nil {
		return errors.Trace(err)
	}
	for _, repo := range repos {
		repo.ns = namespaceSet[repo.NamespaceID]
	}
	return nil
}

func FillUsers(tx sqlx.Ext, repos ...*Repository) error {
	ownerIDs := make([]int64, 0, len(repos))
	for _, repo := range repos {
		ownerIDs = append(ownerIDs, repo.OwnerID)
	}
	userSet, err := user.MapAllUsersByIDs(tx, ownerIDs...)
	if err != nil {
		return errors.Trace(err)
	}
	for _, repo := range repos {
		repo.owner = userSet[repo.OwnerID]
	}
	return nil
}

type model struct {
	*base.Model
	src sqlx.Ext
}

func New(src sqlx.Ext) *model {
	return &model{
		src:   src,
		Model: base.NewModel(src, tableName, nil),
	}
}
