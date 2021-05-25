package db

func init() {
	DB.Register(&Org{})
}

// OrgUser
// @gengo:tablemodel
// @def primary ID
// @def index i_org_id OrgID
// @def index i_user_id UserID
// organization
type OrgUser struct {
	WithPrimaryID

	// @rel Org.OrgID
	// 关联组织
	OrgID uint64 `db:"f_org_id"`
	// @rel User.ID
	// 关联用户
	UserID uint64 `db:"f_user_id"`
}
