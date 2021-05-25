package querysql

import (
	"fmt"
	"strings"
)

func ParseJoinOn(dsl string, r PkgPathResolver) (JoinOn, error) {
	parts := strings.Split(dsl, "=")
	if len(parts) != 2 {
		return JoinOn{}, fmt.Errorf("invalid joind %s, need `<PkgName>.<Export>#<FieldName> = <PkgName>.<Export2>#<FieldName>`", dsl)
	}
	return JoinOn{}, nil
}

type JoinOn [2]*ModelField

type PkgPathResolver = func(localName string) string

func ParseModelField(dsl string, r PkgPathResolver) (*ModelField, error) {
	modelField := strings.Split(dsl, "#")

	if len(modelField) != 2 {
		return nil, fmt.Errorf("invalid where def %s, need `<PkgName>.<Export>#<FieldName>`", dsl)
	}

	name, localName := exportedName(modelField[0])

	pkgPath := r(localName)

	if pkgPath == "" {
		return nil, fmt.Errorf("invalid where def %s, need `<PkgName>.<Export>#<FieldName>`", dsl)
	}

	wm := &ModelField{
		PkgPath:   pkgPath,
		Name:      name,
		FieldName: modelField[1],
	}

	return wm, nil
}

type ModelField struct {
	PkgPath   string
	Name      string
	FieldName string
}

func exportedName(s string) (name string, pkgNameOrPkgPath string) {
	parts := strings.Split(s, ".")
	n := len(parts)
	if n == 0 {
		return name, ""
	}
	return parts[n-1], strings.Join(parts[0:n-1], ".")
}
