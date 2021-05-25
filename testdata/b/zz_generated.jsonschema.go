/*
Package b GENERATED BY gengo:jsonschema
DON'T EDIT THIS FILE
*/
package b

import (
	github_com_go_courier_schema_pkg_jsonschema "github.com/go-courier/schema/pkg/jsonschema"
	github_com_go_courier_schema_pkg_jsonschema_extractors "github.com/go-courier/schema/pkg/jsonschema/extractors"
)

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/go-courier/schema/testdata/b.PullPolicy", new(PullPolicy))
}

func (PullPolicy) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"string",
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Enum: []interface{}{
					"Always",
					"IfNotPresent",
					"Never",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-enum-labels": []string{
					"Always",
					"if not preset",
					"Never",
				},
				"x-go-vendor-type": "github.com/go-courier/schema/testdata/b.PullPolicy",
			},
		},
	})
}
