package namespace

import "github.com/growerlab/growerlab/src/common/errors"

type Namespace struct {
	ID      int64  `db:"id"`
	Path    string `db:"path"`
	OwnerID int64  `db:"owner_id"`
	Type    int    `db:"type"`
}

func (n *Namespace) TypeLabel() string {
	switch NamespaceType(n.Type) {
	case TypeUser:
		return TypeUserLabel
	case TypeOrg:
		return TypeOrgLabel
	default:
		panic(errors.Errorf("invalid namespace type '%d'", n.Type))
	}
}
