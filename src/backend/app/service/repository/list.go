package repository

import (
	"github.com/gin-gonic/gin"
	namespaceModel "github.com/growerlab/growerlab/src/backend/app/model/namespace"
	repositoryModel "github.com/growerlab/growerlab/src/backend/app/model/repository"
	"github.com/growerlab/growerlab/src/backend/app/service/common/session"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/permission"
)

func ListRepositories(c *gin.Context, namespace string) ([]*repositoryModel.Repository, error) {
	currentUserNSID := session.New(c).UserNamespace()

	ns, err := namespaceModel.GetNamespaceByPath(db.DB, namespace)
	if err != nil {
		return nil, err
	}
	if ns == nil {
		return nil, errors.NotFoundError(errors.Namespace)
	}

	repositories, err := repositoryModel.New(db.DB).ListRepositoriesByNamespace(ns.ID)
	if err != nil {
		return nil, err
	}

	var result []*repositoryModel.Repository
	for _, repo := range repositories {
		err := permission.CheckViewRepository(currentUserNSID, repo.ID)
		if err == nil {
			result = append(result, repo)
		}
	}

	return result, nil
}
