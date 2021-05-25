package querysql

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestParseModelField(t *testing.T) {
	p, err := ParseModelField("db.User#ID", func(localName string) string {
		return "xxx/db"
	})
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(p).To(Equal(&ModelField{
		PkgPath:   "xxx/db",
		Name:      "User",
		FieldName: "ID",
	}))
}
