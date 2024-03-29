package userdomain

const (
	AnonymousVisitor = -1 // 特殊用户id：访客
)

// 用户域
const (
	TypeSuperAdmin      = 2000 // 超级管理员
	TypePerson          = 2001 // 个人
	TypeRepositoryOwner = 2002 // 仓库创建者
	TypeVisitor         = 2003 // 每个人（含访客）
)
