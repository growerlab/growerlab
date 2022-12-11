package repository

import (
	"log"

	namespaceModel "github.com/growerlab/growerlab/src/backend/app/model/namespace"
	repositoryModel "github.com/growerlab/growerlab/src/backend/app/model/repository"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/permission"
	"github.com/samber/lo"
)

func (g *Take) List() (*ListResponse, error) {
	ns, err := namespaceModel.GetNamespaceByPath(db.DB, g.namespace)
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

	repos := lo.Filter(repositories, func(item *repositoryModel.Repository, _ int) bool {
		err = permission.CheckViewRepository(g.currentUserID, item.ID)
		if err != nil {
			log.Printf("can't view the repo '%s'", item.PathGroup())
		}
		return err == nil
	})

	err = repositoryModel.FillNamespaces(db.DB, repos...)
	if err != nil {
		return nil, err
	}
	err = repositoryModel.FillUsers(db.DB, repos...)
	if err != nil {
		return nil, err
	}

	var result = make([]*Entity, 0, len(repos))
	for _, repo := range repos {
		result = append(result, BuildRepositryEntity(repo))
	}

	return &ListResponse{Repositories: result}, nil
}
