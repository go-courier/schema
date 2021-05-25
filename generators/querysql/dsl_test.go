package querysql

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestParseWhereModel(t *testing.T) {
	p, err := ParseWhereModel("db.User#ID", func(localName string) string {
		return "xxx/db"
	})
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(p).To(Equal(&WhereModel{
		PkgPath:   "xxx/db",
		Name:      "User",
		FieldName: "ID",
	}))
}
