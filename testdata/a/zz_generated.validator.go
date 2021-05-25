/*
Package a GENERATED BY gengo:validator
DON'T EDIT THIS FILE
*/
package a

import (
	github_com_go_courier_schema_pkg_validator "github.com/go-courier/schema/pkg/validator"
	github_com_go_courier_schema_testdata_b "github.com/go-courier/schema/testdata/b"
)

func validateStruct(v *Struct) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateStructFieldInt(&v.Int), "int")
	errSet.AddErr(validateStructFieldName(&v.Name), "name")
	errSet.AddErr(validateStructFieldID(&v.ID), "id")
	errSet.AddErr(validateStructFieldPullPolicy(&v.PullPolicy), "pullPolicy")
	errSet.AddErr(validateStructFieldProtocol(&v.Protocol), "protocol")
	errSet.AddErr(validateStructFieldSlice(&v.Slice), "slice")
	errSet.AddErr(validateStructFieldMap(&v.Map), "map")
	errSet.AddErr(validateStructFieldC(&v.C), "")
	errSet.AddErr(validateStructFieldChildren(&v.Children), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateStructFieldInt(v *int) error {
	vv := *v
	if vv == 0 {
		vv = 7
		*v = vv
	}
	return (&(github_com_go_courier_schema_pkg_validator.IntValidator{
		BitSize:          32,
		Minimum:          func(v int64) *int64 { return &v }(-2147483648),
		Maximum:          func(v int64) *int64 { return &v }(1024),
		ExclusiveMaximum: true,
		ExclusiveMinimum: true,
	})).Validate(int64(vv))

}
func validateStructFieldName(v **string) error {
	vv := *v
	if vv == nil {
		return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
	}
	return validateStructFieldNamePtr(vv)

}
func validateStructFieldNamePtr(v *string) error {
	vv := *v
	if vv == "" {
		return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
	}
	return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
		MinLength: 2,
	})).Validate(string(vv))

}
func validateStructFieldID(v ***string) error {
	vv := *v
	if vv == nil {
		var vvv *string
		vv = &vvv
		*v = vv
	}
	return validateStructFieldIDPtr(vv)

}
func validateStructFieldIDPtr(v **string) error {
	vv := *v
	if vv == nil {
		var vvv string
		vv = &vvv
		*v = vv
	}
	return validateStructFieldIDPtrPtr(vv)

}
func validateStructFieldIDPtrPtr(v *string) error {
	vv := *v
	if vv == "" {
		vv = "11"
		*v = vv
	}
	return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
		Pattern: "\\d+",
	})).Validate(string(vv))

}
func validateStructFieldPullPolicy(v *github_com_go_courier_schema_testdata_b.PullPolicy) error {
	vv := *v
	if vv == "" {
		return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
	}
	return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
		Enums: []string{
			"Always",
		},
	})).Validate(string(vv))

}
func validateStructFieldProtocol(v *Protocol) error {
	vv, _ := v.MarshalText()
	if vv != nil {
		return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
			Enums: []string{
				"HTTP",
			},
		})).Validate(string(vv))
	}
	return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})

}
func validateStructFieldSlice(v *[]int) error {
	list := *v
	n := len(list)
	if n < 1 {
		return &(github_com_go_courier_schema_pkg_validator.OutOfRangeError{
			Target:  "slice length",
			Current: n,
			Minimum: 1,
		})
	}
	if n > 30 {
		return &(github_com_go_courier_schema_pkg_validator.OutOfRangeError{
			Target:  "slice length",
			Current: n,
			Maximum: 30,
		})
	}
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for i := 0; i < n; i++ {
		errSet.AddErr(validateStructFieldSliceElem(&list[i]), i)
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateStructFieldSliceElem(v *int) error {
	vv := *v
	if vv == 0 {
		return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
	}
	return (&(github_com_go_courier_schema_pkg_validator.IntValidator{
		BitSize: 32,
		Minimum: func(v int64) *int64 { return &v }(0),
		Maximum: func(v int64) *int64 { return &v }(10),
	})).Validate(int64(vv))

}
func validateStructFieldMap(v *map[string]map[string]struct {
	ID int `json:"id" validate:"@int[0,10]"`
}) error {
	m := *v
	if n := len(m); n > 3 {
		return &(github_com_go_courier_schema_pkg_validator.OutOfRangeError{
			Target:  "slice length",
			Current: n,
			Maximum: 3,
		})
	}
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for k := range m {

		value := m[k]
		if e := validateStructFieldMapElem(&value); e != nil {
			errSet.AddErr(e, k)
		}
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateStructFieldMapElem(v *map[string]struct {
	ID int `json:"id" validate:"@int[0,10]"`
}) error {
	m := *v
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for k := range m {

		value := m[k]
		if e := validateStructFieldMapElemElem(&value); e != nil {
			errSet.AddErr(e, k)
		}
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateStructFieldMapElemElem(v *struct {
	ID int `json:"id" validate:"@int[0,10]"`
}) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateStructFieldMapElemElemFieldID(&v.ID), "id")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateStructFieldMapElemElemFieldID(v *int) error {
	vv := *v
	if vv == 0 {
		return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
	}
	return (&(github_com_go_courier_schema_pkg_validator.IntValidator{
		BitSize: 32,
		Minimum: func(v int64) *int64 { return &v }(0),
		Maximum: func(v int64) *int64 { return &v }(10),
	})).Validate(int64(vv))

}
func validateStructFieldC(v **Struct) error {
	vv := *v
	if vv == nil {
		return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
	}
	return validateStructFieldCPtr(vv)

}
func validateStructFieldCPtr(v *Struct) error {
	return validateStruct(v)
}
func validateStructFieldChildren(v *[]*Struct) error {
	list := *v
	n := len(list)
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for i := 0; i < n; i++ {
		errSet.AddErr(validateStructFieldChildrenElem(&list[i]), i)
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateStructFieldChildrenElem(v **Struct) error {
	vv := *v
	if vv == nil {
		return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
	}
	return validateStructFieldChildrenElemPtr(vv)

}
func validateStructFieldChildrenElemPtr(v *Struct) error {
	return validateStruct(v)
}
func (v *Struct) Validate() error {
	return validateStruct(v)
}
