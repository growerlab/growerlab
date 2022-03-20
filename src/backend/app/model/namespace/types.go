package namespace

type Namespace struct {
	ID      int64  `db:"id"`
	Path    string `db:"path"`
	OwnerID int64  `db:"owner_id"`
	Type    int    `db:"type"`
}
