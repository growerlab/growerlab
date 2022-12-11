package userdomain

import (
	"github.com/growerlab/growerlab/src/backend/app/model/repository"
)

type RepositoryOwner struct {
}

func (s *RepositoryOwner) Type() int {
	return TypeRepositoryOwner
}

func (s *RepositoryOwner) TypeLabel() string {
	return "repository_owner"
}

func (s *RepositoryOwner) Validate(ud *UserDomain) error {
	return nil
}

func (s *RepositoryOwner) Eval(args Evaluable) ([]int64, error) {
	repoID := args.Context().Param1
	repo, err := repository.New(args.DB().Src).GetRepository(repoID)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return []int64{}, nil
	}

	ownerID := repo.Owner().ID
	return []int64{ownerID}, nil
}
