package querysql

import (
	"fmt"
	"strings"
)

type PkgPathResolver = func(localName string) string

// ParseWhereModel
// db.User#CreatedAt
func ParseWhereModel(dsl string, r PkgPathResolver) (*WhereModel, error) {
	modelField := strings.Split(dsl, "#")

	if len(modelField) != 2 {
		return nil, fmt.Errorf("invalid where def %s, need `<PkgName>.<Export>#<FieldName>`", dsl)
	}

	name, localName := export(modelField[0])

	pkgPath := r(localName)

	if pkgPath == "" {
		return nil, fmt.Errorf("invalid where def %s, need `<PkgName>.<Export>#<FieldName>`", dsl)
	}

	wm := &WhereModel{
		PkgPath: pkgPath,
		Name: name,
		FieldName:      modelField[1],
	}

	return wm, nil
}

type WhereModel struct {
	PkgPath string
	Name string
	FieldName      string
}

func export(s string) (name string, pkgNameOrPkgPath string) {
	parts := strings.Split(s, ".")
	n := len(parts)
	if n == 0 {
		return name, ""
	}
	return parts[n-1], strings.Join(parts[0:n-1], ".")
}
