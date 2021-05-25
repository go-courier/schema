package db

func init() {
	DB.Register(&Org{})
}

// Org
// @gengo:tablemodel
// @def unique_index i_org_id OrgID
// organization
type Org struct {
	OrgID uint64 `db:"f_org_id"`
	Name  string `db:"f_name,default=''"`
}
