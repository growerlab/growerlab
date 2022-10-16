package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/backend/app/common/permission"
	namespaceModel "github.com/growerlab/growerlab/src/backend/app/model/namespace"
	repositoryModel "github.com/growerlab/growerlab/src/backend/app/model/repository"
	"github.com/growerlab/growerlab/src/backend/app/service/common/session"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
)

func GetRepository(c *gin.Context, namespace, path string) (*repositoryModel.Repository, error) {
	if namespace == "" {
		return nil, errors.InvalidParameterError(errors.Namespace, errors.Path, errors.Empty)
	}
	if path == "" {
		return nil, errors.InvalidParameterError(errors.Repository, errors.Path, errors.Empty)
	}

	currentUserNSID := session.New(c).UserNamespace()

	ns, err := namespaceModel.GetNamespaceByPath(db.DB, namespace)
	if err != nil {
		return nil, err
	}
	if ns == nil {
		return nil, errors.NotFoundError(errors.Namespace)
	}

	repo, err := repositoryModel.GetRepositoryByNsWithPath(db.DB, ns.ID, path)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, errors.NotFoundError(errors.Repository)
	}

	err = permission.CheckViewRepository(currentUserNSID, repo.ID)
	if err != nil {
		return nil, err
	}
	return repo, err
}
