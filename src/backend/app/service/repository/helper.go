package repository

import (
	"time"

	"github.com/growerlab/growerlab/src/backend/app/model/repository"
	"github.com/growerlab/growerlab/src/backend/app/utils/uuid"
)

func BuildNewRepository(
	userID int64,
	nsID int64,
	req *CreateParams,
) (repo *repository.Repository) {
	repo = &repository.Repository{
		NamespaceID: nsID,
		UUID:        uuid.UUIDv16(),
		Path:        req.Name,
		Name:        req.Name,
		OwnerID:     userID,
		Description: req.Description,
		CreatedAt:   time.Now().Unix(),
		Public:      req.Public,
	}
	return repo
}

func BuildRepositryEntity(
	repo *repository.Repository) *Entity {

	return &Entity{
		UUID:        repo.UUID,
		Name:        repo.Name,
		Path:        repo.Path,
		Description: repo.Description,
		CreatedAt:   repo.CreatedAt,
		Public:      repo.Public,
		Namespace: &NamespaceEntity{
			Path: repo.Namespace().Path,
			Type: repo.Namespace().TypeLabel(),
		},
		Owner: &UserEntity{
			Username:    repo.Owner().Username,
			Name:        repo.Owner().Name,
			PublicEmail: repo.Owner().Email,
		},
	}
}
