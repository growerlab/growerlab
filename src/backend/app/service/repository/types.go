package repository

type ListResponse struct {
	Repositories []*Entity `json:"repositories"`
}

type Entity struct {
	UUID          string `json:"uuid"` // 全站唯一ID（fork时用到）
	Name          string `json:"name"` // 目前与path字段相同
	Path          string `json:"path"` // 在namespace中唯一
	Description   string `json:"description"`
	CreatedAt     int64  `json:"created_at"`
	Public        bool   `json:"public"`         // 公有
	LastPushAt    int64  `json:"last_push_at"`   // 最后的推送时间
	DefaultBranch string `json:"default_branch"` // 默认分支

	Namespace *NamespaceEntity `json:"namespace"`
	Owner     *UserEntity      `json:"owner"`
}

type NamespaceEntity struct {
	Path string `json:"path"`
	Type string `json:"type"` // 个人/团队
}

type UserEntity struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	PublicEmail string `json:"public_email"`
}
