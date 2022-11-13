package repository

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/backend/app/common/git"
	"github.com/growerlab/growerlab/src/backend/app/model/namespace"
	"github.com/growerlab/growerlab/src/backend/app/model/repository"
	"github.com/growerlab/growerlab/src/backend/app/model/server"
	"github.com/growerlab/growerlab/src/backend/app/model/user"
	"github.com/growerlab/growerlab/src/backend/app/service/common/session"
	"github.com/growerlab/growerlab/src/backend/app/utils/regex"
	"github.com/growerlab/growerlab/src/backend/app/utils/uuid"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/path"
	"github.com/jmoiron/sqlx"
)

type NewRepositoryPayload struct {
	Namespace   string `json:"namespace"` // 命名空间的路径（这里要考虑某个人在组织下创建项目）
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	Description string `json:"description"`
}

func CreateRepository(c *gin.Context, req *NewRepositoryPayload) error {
	currentUser := session.New(c).User()
	if currentUser == nil {
		return errors.UnauthorizedError()
	}
	return DoCreateRepository(currentUser, req)
}

func DoCreateRepository(currentUser *user.User, req *NewRepositoryPayload) error {
	var ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	err := db.Transact(func(tx sqlx.Ext) error {
		ns, err := validateAndPrepare(tx, currentUser.ID, req)
		if err != nil {
			return errors.Trace(err)
		}

		srv, err := server.RandNormalServer(tx)
		if err != nil {
			return errors.Trace(err)
		}

		repo := buildRepository(currentUser, ns, req, srv)
		if err != nil {
			return errors.Trace(err)
		}

		err = repository.New(tx).AddRepository(repo)
		if err != nil {
			return errors.Trace(err)
		}

		// 真正创建仓库
		err = git.New(ctx, repo.PathGroup()).Create()
		if err != nil {
			return errors.Trace(err)
		}
		return nil
	})
	return err
}

func buildRepository(
	currentUser *user.User,
	ns *namespace.Namespace,
	req *NewRepositoryPayload,
	srv *server.Server,
) (repo *repository.Repository) {
	serverPath := path.GetRelativeRepositoryPath(filepath.Join(ns.Path, req.Name))
	repo = &repository.Repository{
		NamespaceID: ns.ID,
		UUID:        uuid.UUIDv16(),
		Path:        req.Name,
		Name:        req.Name,
		OwnerID:     currentUser.ID,
		Description: req.Description,
		CreatedAt:   time.Now().Unix(),
		ServerID:    srv.ID,
		ServerPath:  serverPath,
		Public:      req.Public,
	}
	return repo
}

// validate
//
//	req.NamespacePath  TODO 这里暂时只验证namespace的owner_id 是否为用户，未来应该验证组织权限（比如是否可以选择这个组织创建仓库）
//	req.Name 名称是否合法、是否重名
func validateAndPrepare(src sqlx.Ext, userID int64, req *NewRepositoryPayload) (ns *namespace.Namespace, err error) {
	req.Namespace = strings.TrimSpace(req.Namespace)
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Namespace) == 0 {
		err = errors.InvalidParameterError(errors.Namespace, errors.Path, errors.Invalid)
		return
	}

	if len(req.Name) == 0 {
		err = errors.InvalidParameterError(errors.Repository, errors.Name, errors.Invalid)
		return
	}

	ns, err = namespace.GetNamespaceByPath(src, req.Namespace)
	if err != nil {
		return nil, err
	}

	// TODO 未来应该验证权限(例如是否有权限在组织中创建权限)
	if ns.OwnerID != userID {
		return nil, errors.AccessDenied(errors.User, errors.NotEqual)
	}

	// 验证仓库名是否合法
	// 1. 不能允许以 .. 符号开始（安全问题）
	// 2. 其余以 regex.RepositoryNameRegex 规则为准
	if strings.Index(req.Name, "..") == 0 {
		return nil, errors.InvalidParameterError(errors.Repository, errors.Name, errors.Invalid)
	}
	if !regex.Match(req.Name, regex.RepositoryNameRegex) {
		return nil, errors.InvalidParameterError(errors.Repository, errors.Name, errors.Invalid)
	}

	// 验证仓库名在repository.namespace中是否已存在
	exist, err := repository.New(src).NameExists(ns.ID, req.Name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.AlreadyExistsError(errors.Repository, errors.AlreadyExists)
	}
	return ns, nil
}
