package a

import (
	github_com_go_courier_schema_pkg_validator "github.com/go-courier/schema/pkg/validator"
)

func (v *Struct) Validate() error {
	validate := func(v *Struct) error {
		errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()

		{
			var fe error
			if fv := int64(v.Int); fv != 0 {
				fe = (&(github_com_go_courier_schema_pkg_validator.IntValidator{})).Validate(fv)
			} else {
				fe = github_com_go_courier_schema_pkg_validator.MissingRequired{}
			}

			if fe != nil {
				errSet.AddErr(fe, "int")
			}
		}
		{
			var fe error
			if v.Name != nil {
				if fv := string(*v.Name); fv != "" {
					fe = (&(github_com_go_courier_schema_pkg_validator.StringValidator{})).Validate(fv)
				} else {
					fe = github_com_go_courier_schema_pkg_validator.MissingRequired{}
				}
			} else {
				fe = github_com_go_courier_schema_pkg_validator.MissingRequired{}
			}

			if fe != nil {
				errSet.AddErr(fe, "name")
			}
		}
		{
			var fe error
			if v.ID != nil && *v.ID != nil {
				if fv := string(**v.ID); fv != "" {
					fe = (&(github_com_go_courier_schema_pkg_validator.StringValidator{})).Validate(fv)
				}
			}
			if fe != nil {
				errSet.AddErr(fe, "id")
			}
		}
		{
			var fe error
			if fv := string(v.PullPolicy); fv != "" {
				fe = (&(github_com_go_courier_schema_pkg_validator.StringValidator{})).Validate(fv)
			} else {
				fe = github_com_go_courier_schema_pkg_validator.MissingRequired{}
			}

			if fe != nil {
				errSet.AddErr(fe, "pullPolicy")
			}
		}
		{
			var fe error
			if fv := github_com_go_courier_schema_pkg_validator.ToMarshalledText(v.Protocol); fv != "" {
				fe = (&(github_com_go_courier_schema_pkg_validator.IntValidator{})).Validate(fv)
			} else {
				fe = github_com_go_courier_schema_pkg_validator.MissingRequired{}
			}

			if fe != nil {
				errSet.AddErr(fe, "protocol")
			}
		}

		return errSet.Err()
	}
	return validate(v)
}
