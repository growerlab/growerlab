package permission

type Permission struct {
	ID              int64  `db:"id"`
	NamespaceID     int64  `db:"namespace_id"`
	Code            int    `db:"code"`
	ContextType     int    `db:"context_type"`
	ContextParam1   int64  `db:"context_param_1"`
	ContextParam2   int64  `db:"context_param_2"`
	UserDomainType  int    `db:"user_domain_type"`
	UserDomainParam int64  `db:"user_domain_param"`
	CreatedAt       int64  `db:"created_at"`
	DeletedAt       *int64 `db:"deleted_at"`
}
