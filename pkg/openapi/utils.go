package openapi

import "github.com/go-courier/schema/pkg/jsonschema"

type OperationGetter interface {
	OpenAPIOperation(ref func(t string) jsonschema.Refer) *Operation
}
