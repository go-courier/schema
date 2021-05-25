package db

import (
	"github.com/go-courier/schema/testdata/a"
	"github.com/go-courier/sqlx/v2/datatypes"
)

func init() {
	DB.Register(&User{})
}

// User
// @gengo:tablemodel
// @def unique_index i_user_id UserID
// @def index i_nickname/BTREE Nickname
// @def index i_username Username
// @def index i_geom/SPATIAL Geom
// @def unique_index i_name Name
type User struct {
	UserID uint64 `db:"f_user_id"`
	// 姓名
	Name     string     `db:"f_name,default=''"`
	Username string     `db:"f_username,default=''"`
	Nickname string     `db:"f_nickname,default=''"`
	Gender   a.Protocol `db:"f_gender,default='0'"`
	Boolean  bool       `db:"f_boolean,default=false"`
	Geom     GeomString `db:"f_geom"`

	CreatedAt datatypes.Timestamp `db:"f_created_at,default='0'"`
	UpdatedAt datatypes.Timestamp `db:"f_updated_at,default='0'"`
	DeletedAt datatypes.Timestamp `db:"f_deleted_at,default='0'"`
}
