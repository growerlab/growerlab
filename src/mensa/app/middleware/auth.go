package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/growerlab/growerlab/src/backend/app/common/permission"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/mensa/app/common"
	"github.com/growerlab/growerlab/src/mensa/app/service"
)

var (
	ErrUnauthorized = MiddlewareError("Unauthorized")
)

// Authenticate 鉴权
func Authenticate(ctx *common.Context) (httpCode int, appendText string, err error) {
	httpCode = http.StatusOK
	noAuth := os.Getenv("NOAUTH")
	if len(noAuth) > 0 {
		return
	}

	if err = checkPermission(ctx); err != nil {
		switch errors.Cause(err) {
		case service.NotFoundRepoError:
			httpCode = http.StatusNotFound
		default:
			httpCode = http.StatusUnauthorized
		}
		appendText = err.Error()
		log.Printf("%s, err: unauthorized '%v'\n", ctx.Desc(), err)
		return
	}
	return
}

// 检查是否有读取、推送权限
//
//	公共项目：可读、只有项目成员可写
//	私有项目：项目成员可读/写
func checkPermission(ctx *common.Context) error {
	repo, err := service.GetRepository(ctx.RepoOwner, ctx.RepoName)
	if err != nil {
		return err
	}

	if repo.IsPublic() {
		if ctx.IsReadAction() {
			return nil
		} else {
			var nsID int64
			var err error
			if ctx.Operator == nil {
				return errors.WithStack(ErrUnauthorized)
			}
			if !ctx.Operator.IsEmptyUser() {
				nsID, err = service.GetNamespaceByOperator(ctx.Operator)
				if err != nil {
					return err
				}
			}
			return permission.CheckPushRepository(nsID, repo.ID)
		}
	} else {
		if ctx.Operator == nil {
			return errors.WithStack(ErrUnauthorized)
		}

		nsID, err := service.GetNamespaceByOperator(ctx.Operator)
		if err != nil {
			return err
		}

		if ctx.IsReadAction() {
			return permission.CheckCloneRepository(&nsID, repo.ID)
		} else {
			return permission.CheckPushRepository(nsID, repo.ID)
		}
	}
}
