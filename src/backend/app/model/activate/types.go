package activate

const CodeMaxLen = 16

type ActivationCode struct {
	ID        int64  `db:"id"`
	UserID    int64  `db:"user_id"`
	Code      string `db:"code"`
	CreatedAt int64  `db:"created_at"`
	UsedAt    *int64 `db:"used_at"`
	ExpiredAt int64  `db:"expired_at"`
}
