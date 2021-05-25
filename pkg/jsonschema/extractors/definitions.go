package extractors

import "github.com/go-courier/schema/pkg/jsonschema"

type OpenAPISchemaGetter interface {
	OpenAPISchema(ref func(t string) jsonschema.Refer) *jsonschema.Schema
}

type OpenAPISchemaTypeGetter interface {
	OpenAPISchemaType() []string
}

type OpenAPISchemaFormatGetter interface {
	OpenAPISchemaFormat() string
}

type Definitions map[string]OpenAPISchemaGetter

func (Definitions) Ref(key string) jsonschema.Refer {
	return jsonschema.RefSchema("#/definitions/" + key)
}

var definitions = Definitions{}

func Get(key string) (OpenAPISchemaGetter, bool) {
	g, ok := definitions[key]
	return g, ok
}

func Register(key string, getter OpenAPISchemaGetter) {
	definitions[key] = getter
}
