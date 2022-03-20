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
	result := make([]int64, 0)
	repoID := args.Context().Param1
	repo, err := repository.GetRepository(args.DB().Src, repoID)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return result, nil
	}

	ownerNamespaceID := repo.Owner().NamespaceID
	return []int64{ownerNamespaceID}, nil
}
