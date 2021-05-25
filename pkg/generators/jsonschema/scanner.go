package jsonschema

import (
	"context"
	"go/token"
	"go/types"
	"strings"

	"github.com/go-courier/gengo/pkg/gengo"
	"github.com/go-courier/schema/pkg/generators/enum"
	"github.com/go-courier/schema/pkg/jsonschema"
	"github.com/go-courier/schema/pkg/jsonschema/extractors"
)

type Scanner struct {
	enumTypes         enum.EnumTypes
	definitions       map[string]*jsonschema.Schema
	vendorDefinitions map[string]*jsonschema.Schema
}

func (g *Scanner) Init(c *gengo.Context) {
	g.enumTypes = enum.EnumTypes{}
	g.enumTypes.Walk(c, c.Package.Pkg().Path())
	g.definitions = map[string]*jsonschema.Schema{}
}

func (g *Scanner) SchemaFromType(ctx context.Context, pos token.Pos, tpe types.Type, def bool) (s *jsonschema.Schema) {
	defer func() {
		if s != nil {
			if def {
				if e, ok := g.enumTypes.ResolveEnumType(tpe); ok {
					values := make([]interface{}, len(e.Constants))
					labels := make([]string, len(e.Constants))

					for i := range e.Constants {
						values[i] = e.Value(e.Constants[i])
						labels[i] = e.Label(e.Constants[i])
					}

					if len(values) > 0 {
						if _, ok := values[0].(string); ok {
							s.Type = jsonschema.StringOrArray{"string"}
						}
					}

					s.Enum = values
					s.AddExtension(jsonschema.XEnumLabels, labels)
				}
			}

			if docGetter := extractors.DocGetterFromContext(ctx); docGetter != nil {
				_, lines := docGetter.Doc(pos)
				s.Description = strings.Join(lines, "\n")
			}
		}
	}()

	return extractors.OpenAPISchemaGetterFromGoType(ctx, tpe, def).OpenAPISchema(func(t string) jsonschema.Refer {
		return extractors.TypeName(t)
	})
}
