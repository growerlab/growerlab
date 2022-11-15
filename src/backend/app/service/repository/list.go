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

func ListRepositories(c *gin.Context, namespace string) ([]*RepositoryEntity, error) {
	currentUserID := session.New(c).UserID()

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

	var repos []*repositoryModel.Repository
	for _, repo := range repositories {
		err := permission.CheckViewRepository(currentUserID, repo.ID)
		if err == nil {
			repos = append(repos, repo)
		}
	}

	err = repositoryModel.FillNamespaces(db.DB, repos...)
	if err != nil {
		return nil, err
	}
	err = repositoryModel.FillUsers(db.DB, repos...)
	if err != nil {
		return nil, err
	}

	var result = make([]*RepositoryEntity, 0, len(repos))
	for _, repo := range repos {
		result = append(result, BuildRepositryEntity(repo))
	}
	return result, nil
}
