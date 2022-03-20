package userdomain

import (
	"github.com/growerlab/growerlab/src/backend/app/model/user"
)

type SuperAdmin struct {
}

func (s *SuperAdmin) Type() int {
	return TypeSuperAdmin
}

func (s *SuperAdmin) TypeLabel() string {
	return "super_admin"
}

func (s *SuperAdmin) Validate(ud *UserDomain) error {
	return nil
}

func (s *SuperAdmin) Eval(args Evaluable) ([]int64, error) {
	admins, err := user.ListAdminUsers(args.DB().Src)
	if err != nil {
		return nil, err
	}
	namespaceIds := make([]int64, len(admins))
	for i := range admins {
		namespaceIds[i] = admins[i].NamespaceID
	}
	return namespaceIds, nil
}
