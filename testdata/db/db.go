package db

import "github.com/go-courier/sqlx/v2"

var DB = sqlx.NewDatabase("test")

type WithPrimaryID struct {
	ID uint64 `db:"f_id,autoincrement" json:"-"`
}
