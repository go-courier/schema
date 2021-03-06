/*
Package routes GENERATED BY gengo:jsonschema
DON'T EDIT THIS FILE
*/
package routes

import (
	mime_multipart "mime/multipart"

	github_com_go_courier_httptransport_httpx "github.com/go-courier/httptransport/httpx"
	github_com_go_courier_schema_pkg_jsonschema "github.com/go-courier/schema/pkg/jsonschema"
	github_com_go_courier_schema_pkg_jsonschema_extractors "github.com/go-courier/schema/pkg/jsonschema/extractors"
)

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/go-courier/schema/testdata/routes.Account", new(Account))
}

func (Account) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
			Properties: map[string]*github_com_go_courier_schema_pkg_jsonschema.Schema{
				"UserID": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "UserID",
						},
					},
				}),
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Required: []string{
					"UserID",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/go-courier/schema/testdata/routes.Account",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/go-courier/schema/testdata/routes.Data", new(Data))
}

func (Data) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
			Properties: map[string]*github_com_go_courier_schema_pkg_jsonschema.Schema{
				"data": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						AllOf: []*github_com_go_courier_schema_pkg_jsonschema.Schema{
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
									Nullable: true,
								},
								Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
									Refer: ref("github.com/go-courier/schema/testdata/routes.Data"),
								},
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-go-star-level": 1,
									},
								},
							}),
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-go-field-name": "Data",
									},
								},
							}),
						},
					},
				}),
				"id": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "ID",
						},
					},
				}),
				"label": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Label",
						},
					},
				}),
				"protocol": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						AllOf: []*github_com_go_courier_schema_pkg_jsonschema.Schema{
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
									Refer: ref("github.com/go-courier/schema/testdata/a.Protocol"),
								},
							}),
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-go-field-name": "Protocol",
									},
								},
							}),
						},
					},
				}),
				"ptrString": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
						Nullable: true,
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "PtrString",
							"x-go-star-level": 1,
						},
					},
				}),
				"subData": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						AllOf: []*github_com_go_courier_schema_pkg_jsonschema.Schema{
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
									Nullable: true,
								},
								Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
									Refer: ref("github.com/go-courier/schema/testdata/routes.SubData"),
								},
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-go-star-level": 1,
									},
								},
							}),
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-go-field-name": "SubData",
									},
								},
							}),
						},
					},
				}),
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Required: []string{
					"id",
					"label",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/go-courier/schema/testdata/routes.Data",
			},
		},
	})
}

type shadowedAttachment struct {
	github_com_go_courier_httptransport_httpx.Attachment
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/go-courier/httptransport/httpx.Attachment", new(shadowedAttachment))
}

func (shadowedAttachment) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"string",
			},
			Format: "binary",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/go-courier/httptransport/httpx.Attachment",
			},
		},
	})
}

type shadowedFileHeader struct {
	mime_multipart.FileHeader
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("mime/multipart.FileHeader", new(shadowedFileHeader))
}

func (shadowedFileHeader) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"string",
			},
			Format:      "binary",
			Description: "A FileHeader describes a file part of a multipart request.",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "mime/multipart.FileHeader",
			},
		},
	})
}

type shadowedMethodGet struct {
	github_com_go_courier_httptransport_httpx.MethodGet
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/go-courier/httptransport/httpx.MethodGet", new(shadowedMethodGet))
}

func (shadowedMethodGet) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/go-courier/httptransport/httpx.MethodGet",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/go-courier/schema/testdata/routes.GetByJSON", new(GetByJSON))
}

func (GetByJSON) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			AllOf: []*github_com_go_courier_schema_pkg_jsonschema.Schema{
				&(github_com_go_courier_schema_pkg_jsonschema.Schema{
					Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
						Refer: ref("github.com/go-courier/httptransport/httpx.MethodGet"),
					},
				}),
				&(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"object",
						},
					},
				}),
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/go-courier/schema/testdata/routes.GetByJSON",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/go-courier/schema/testdata/routes.IpInfo", new(IpInfo))
}

func (IpInfo) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
			Properties: map[string]*github_com_go_courier_schema_pkg_jsonschema.Schema{
				"country": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Country",
						},
					},
				}),
				"countryCode": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "CountryCode",
						},
					},
				}),
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Required: []string{
					"country",
					"countryCode",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/go-courier/schema/testdata/routes.IpInfo",
			},
		},
	})
}

type shadowedImagePNG struct {
	github_com_go_courier_httptransport_httpx.ImagePNG
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/go-courier/httptransport/httpx.ImagePNG", new(shadowedImagePNG))
}

func (shadowedImagePNG) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"string",
			},
			Format: "binary",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/go-courier/httptransport/httpx.ImagePNG",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/go-courier/schema/testdata/routes.SubData", new(SubData))
}

func (SubData) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
			Properties: map[string]*github_com_go_courier_schema_pkg_jsonschema.Schema{
				"name": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Name",
						},
					},
				}),
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Required: []string{
					"name",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/go-courier/schema/testdata/routes.SubData",
			},
		},
	})
}
