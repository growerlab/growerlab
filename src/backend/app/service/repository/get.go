package repository

import (
	namespaceModel "github.com/growerlab/growerlab/src/backend/app/model/namespace"
	repositoryModel "github.com/growerlab/growerlab/src/backend/app/model/repository"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/permission"
)

func (g *Repository) Get() (*Entity, error) {
	if g.namespace == "" {
		return nil, errors.InvalidParameterError(errors.Repository, errors.Namespace, errors.Empty)
	}
	if g.repo == nil {
		return nil, errors.InvalidParameterError(errors.Repository, errors.Repo, errors.Empty)
	}

	ns, err := namespaceModel.GetNamespaceByPath(db.DB, g.namespace)
	if err != nil {
		return nil, err
	}
	if ns == nil {
		return nil, errors.NotFoundError(errors.Namespace)
	}

	repo, err := repositoryModel.New(db.DB).GetRepositoryByNsWithPath(ns.ID, *g.repo)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, errors.NotFoundError(errors.Repository)
	}
	err = permission.CheckViewRepository(g.currentUserID, repo.ID)
	if err != nil {
		return nil, err
	}

	err = repositoryModel.FillNamespaces(db.DB, repo)
	if err != nil {
		return nil, err
	}
	err = repositoryModel.FillUsers(db.DB, repo)
	if err != nil {
		return nil, err
	}

	return BuildRepositoryEntity(repo), nil
}
