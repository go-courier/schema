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

// @gengo:querysql
type UserParams struct {
	// @where db.User#Name
	Name RightLikeOrIn `name:"name,omitempty" in:"query"`
	// @where db.User#Username
	Username []string `name:"username,omitempty" in:"query"`
	// @where db.User#Nickname
	Nickname []string `name:"nickname,omitempty" in:"query"`
	// @where db.User#CreatedAt
	CreatedAt DateTimeOrRange `name:"createdAt,omitempty" in:"query"`
}


// @gengo:querysql
// @from db.OrgUser
type OrgUserParams struct {
	// @join db.User#UserID = db.OrgUser#UserID
	UserParams
	// @join db.Org#OrgID = db.OrgUser#OrgID
	OrgParams
}

//func (p *OrgUserParams) ToCondition(d sqlx.TableResolver) builder.SqlCondition {
//	return builder.And(p.UserParams.ToCondition(d), p.OrgParams.ToCondition(d))
//}

type OrgUser struct {
	db.User
	Org db.Org `json:"org"`
}

type OrgUserDataList struct {
	Data  []OrgUser `json:"data"`
	Total int       `json:"total,omitempty"`
}
