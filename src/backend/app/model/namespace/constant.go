package namespace

type NamespaceType int

const (
	TypeUser NamespaceType = 1
	TypeOrg  NamespaceType = 2
)

const (
	TypeUserLabel string = "user"
	TypeOrgLabel  string = "org"
)
