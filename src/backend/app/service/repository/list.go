package repository

import (
	namespaceModel "github.com/growerlab/growerlab/src/backend/app/model/namespace"
	repositoryModel "github.com/growerlab/growerlab/src/backend/app/model/repository"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/permission"
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

	var repos []*repositoryModel.Repository
	for _, repo := range repositories {
		err := permission.CheckViewRepository(g.currentUserID, repo.ID)
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

	var result = make([]*Entity, 0, len(repos))
	for _, repo := range repos {
		result = append(result, BuildRepositryEntity(repo))
	}

	return &ListResponse{Repositories: result}, nil
}
