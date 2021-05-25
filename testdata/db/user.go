package db

import (
	"database/sql/driver"

	"github.com/go-courier/schema/testdata/a"
	"github.com/go-courier/sqlx/v2/datatypes"
)

// +gengo:table
// @def primary ID
// @def index i_nickname/BTREE Nickname
// @def index i_username Username
// @def index i_geom/SPATIAL Geom
// @def unique_index I_name Name
type User struct {
	WithPrimaryID
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

type WithPrimaryID struct {
	ID uint64 `db:"f_id,autoincrement"`
}

type GeomString struct {
	V string
}

func (g GeomString) Value() (driver.Value, error) {
	return g.V, nil
}

func (g *GeomString) Scan(src interface{}) error {
	return nil
}

func (GeomString) DataType(driverName string) string {
	if driverName == "mysql" {
		return "geometry"
	}
	return "geometry(Point)"
}

func (GeomString) ValueEx() string {
	return "ST_GeomFromText(?)"
}
