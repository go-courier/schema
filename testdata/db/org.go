package db

// +gengo:table
// @def primary ID
// organization
type Org struct {
	WithPrimaryID

	Name string `db:"f_name,default=''"`
	// @rel User.ID
	// 关联用户
	// xxxxx
	UserID string `db:"user_id"`
}
