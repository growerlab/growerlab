package permission

import (
	"github.com/growerlab/growerlab/src/backend/app/common/userdomain"
	"github.com/growerlab/growerlab/src/common/context"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/jmoiron/sqlx"
)

var permHub *Hub

func InitPermission() error {
	return InitPermissionHub(db.DB, db.MemDB)
}

func InitPermissionHub(dbSrc sqlx.Queryer, memDB *db.MemDBClient) error {
	permHub = NewPermissionHub(dbSrc, memDB)

	if err := initRules(); err != nil {
		return err
	}
	if err := initUserDomains(); err != nil {
		return err
	}
	if err := initContexts(); err != nil {
		return err
	}
	return nil
}

func initUserDomains() error {
	userDomains := []UserDomainDelegate{
		&userdomain.SuperAdmin{},
		&userdomain.Person{},
		&userdomain.RepositoryOwner{},
		&userdomain.Visitor{},
	}
	return permHub.RegisterUserDomains(userDomains)
}

func initContexts() error {
	contexts := make([]ContextDelegate, 0)
	contexts = append(contexts, &context.Repository{})
	return permHub.RegisterContexts(contexts)
}

func initRules() error {
	rules := []*Rule{
		{
			Code:                  ViewRepository,
			ConstraintUserDomains: []int{userdomain.TypePerson},
			BuiltInUserDomains:    []int{userdomain.TypeRepositoryOwner},
		},
		{
			Code:                  CloneRepository,
			ConstraintUserDomains: []int{userdomain.TypePerson},
			BuiltInUserDomains:    []int{userdomain.TypeRepositoryOwner},
		},
		{
			Code:                  PushRepository,
			ConstraintUserDomains: []int{userdomain.TypePerson},
			BuiltInUserDomains:    []int{userdomain.TypeRepositoryOwner},
		},
	}
	return permHub.RegisterRules(rules)
}
