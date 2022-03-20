package permission

import (
	"github.com/growerlab/growerlab/src/backend/app/common/context"
	"github.com/growerlab/growerlab/src/backend/app/common/userdomain"
)

func NewEvalArgs(ctx *context.Context, ud *userdomain.UserDomain, dbctx *context.DBContext) *EvalArgs {
	return &EvalArgs{
		ctx:   ctx,
		ud:    ud,
		dbctx: dbctx,
	}
}

type EvalArgs struct {
	// 上下文
	ctx *context.Context
	// 大部分情况下，用户域依赖上下文
	ud *userdomain.UserDomain
	// 传入的db上下文
	dbctx *context.DBContext
}

func (e *EvalArgs) UserDomain() *userdomain.UserDomain {
	return e.ud
}

func (e *EvalArgs) Context() *context.Context {
	return e.ctx
}

func (e *EvalArgs) DB() *context.DBContext {
	return e.dbctx
}

type ContextDelegate interface {
	Type() int
	TypeLabel() string
	// Validate 用于新增权限时，对context的参数进行验证，以确保其参数是正确或必填的
	Validate(c *context.Context) error
}

type UserDomainDelegate interface {
	Type() int
	TypeLabel() string
	// Validate 用于新增权限时，对userDomain的参数进行验证，以确保其参数是正确或必填的
	Validate(ud *userdomain.UserDomain) error
	// Eval 根据用户域返回相关的namespace IDs
	Eval(args userdomain.Evaluable) ([]int64, error)
}
