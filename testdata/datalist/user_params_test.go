package datalist

import (
	"context"
	"testing"

	"github.com/go-courier/schema/testdata/db"
)

func BenchmarkUserParams(b *testing.B) {
	u := &OrgUserParams{}
	u.Name = []string{"1", "2", "3"}
	u.Username = []string{"1", "2", "3"}

	e := u.ToCondition(db.DB).Ex(context.Background())
	b.Log(e.Query(), e.Args())

	b.Run("Ex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = u.ToCondition(db.DB).Ex(context.Background())
		}
	})
}
