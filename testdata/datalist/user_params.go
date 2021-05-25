package datalist

import (
	"github.com/go-courier/schema/testdata/db"
)

// @gengo:querysql
type OrgParams struct {
	// @where db.Org#OrgID
	OrgID []uint64 `name:"orgID,omitempty" in:"query"`
}

//func (p *OrgParams) ToCondition(d sqlx.TableResolver) builder.SqlCondition {
//	m := &db.Org{}
//	tUser := d.T(m)
//
//	where := builder.EmptyCond()
//
//	if len(p.OrgID) != 0 {
//		where = builder.And(where, tUser.F(m.FieldKeyOrgID()).In(p.OrgID))
//	}
//
//	return where
//}

// UserParams
// @gengo:querysql
type UserParams struct {
	// @where db.User#Name
	Name RightLikeOrIn `name:"name,omitempty" in:"query"`
	// @where db.User#Username
	Username []string `name:"username,omitempty" in:"query"`
	// @where db.User#Nickname
	Nickname []string `name:"nickname,omitempty" in:"query"`
	// @where db.User#CreatedAt
	UserCreatedAt DateTimeOrRange `name:"userCreatedAt,omitempty" in:"query"`
}

// OrgUserParams
// @gengo:querysql
type OrgUserParams struct {
	// @on db.User#UserID = db.OrgUser#UserID
	UserParams
	// @on db.Org#OrgID = db.OrgUser#OrgID
	OrgParams
}

type OrgUser struct {
	db.User
	Org db.Org `json:"org"`
}

type OrgUserDataList struct {
	Data  []OrgUser `json:"data"`
	Total int       `json:"total,omitempty"`
}
