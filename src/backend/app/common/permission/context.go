package permission

import (
	"github.com/growerlab/growerlab/src/common/context"
)

func RepositoryContext(repositoryID int64) *context.Context {
	return &context.Context{
		Type:   context.TypeRepository,
		Param1: repositoryID,
	}
}
