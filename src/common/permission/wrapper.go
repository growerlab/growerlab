package permission

import (
	"github.com/growerlab/growerlab/src/common/context"
	"github.com/growerlab/growerlab/src/common/userdomain"
)

func CheckViewRepository(userID *int64, repositoryID int64) error {
	c := RepositoryContext(repositoryID)
	return checkPermission(userID, c, ViewRepository)
}

func CheckPushRepository(userID int64, repositoryID int64) error {
	c := RepositoryContext(repositoryID)
	return checkPermission(&userID, c, PushRepository)
}

func CheckCloneRepository(userID *int64, repositoryID int64) error {
	c := RepositoryContext(repositoryID)
	return checkPermission(userID, c, CloneRepository)
}

func checkPermission(userID *int64, ctx *context.Context, code int) error {
	if userID == nil || *userID == 0 {
		userID = new(int64)
		*userID = userdomain.AnonymousVisitor
	}
	return permHub.CheckCache(*userID, ctx, code, true)
}
