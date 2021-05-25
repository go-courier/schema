package datalist

import (
	"context"
	"testing"
	"time"

	"github.com/go-courier/schema/testdata/db"
	"github.com/go-courier/sqlx/v2/datatypes"
)

func BenchmarkUserParams(b *testing.B) {
	u := &UserParams{}
	u.Name = []string{"1", "2", "3"}
	u.Username = []string{"1", "2", "3"}
	u.UserCreatedAt.To = datatypes.Timestamp(time.Now())

	e := u.ToCondition(db.DB).Ex(context.Background())
	b.Log(e.Query(), e.Args())

	b.Run("Ex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = u.ToCondition(db.DB).Ex(context.Background())
		}
	})
}
