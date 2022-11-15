package repository

type RepositoryEntity struct {
	UUID        string `json:"uuid"` // 全站唯一ID（fork时用到）
	Name        string `json:"name"` // 目前与path字段相同
	Path        string `json:"path"` // 在namespace中唯一
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	Public      bool   `json:"public"` // 公有

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
