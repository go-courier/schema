package jsonschema

import (
	"context"
	"fmt"
	"go/types"
	"io"
	"strconv"
	"strings"

	"github.com/go-courier/schema/pkg/generators/enum"
	"github.com/go-courier/schema/pkg/jsonschema/extractors"

	"github.com/go-courier/gengo/pkg/gengo"
	"github.com/go-courier/gengo/pkg/namer"
	gengotypes "github.com/go-courier/gengo/pkg/types"
	"github.com/go-courier/schema/pkg/jsonschema"
)

func init() {
	gengo.Register(&jsonSchemaGen{})
}

type jsonSchemaGen struct {
	imports     namer.ImportTracker
	enumTypes   enum.EnumTypes
	definitions map[string]*jsonschema.Schema
}

func (g *jsonSchemaGen) Name() string {
	return "jsonschema"
}

func (jsonSchemaGen) New() gengo.Generator {
	return &jsonSchemaGen{
		imports:     namer.NewDefaultImportTracker(),
		definitions: map[string]*jsonschema.Schema{},
	}
}

func (g *jsonSchemaGen) Init(c *gengo.Context, w io.Writer) error {
	g.enumTypes = enum.EnumTypes{}
	g.enumTypes.Walk(c, c.Package.Pkg().Path())
	return nil
}

func (g *jsonSchemaGen) Imports(c *gengo.Context) map[string]string {
	return g.imports.Imports()
}

func (g *jsonSchemaGen) GenerateType(c *gengo.Context, named *types.Named, w io.Writer) error {
	sw := gengo.NewSnippetWriter(w, namer.NameSystems{
		"raw": namer.NewRawNamer(c.Package.Pkg().Path(), g.imports),
	})

	a := map[string]interface{}{
		"typeName":    named.Obj().Name(),
		"typePkgPath": named.Obj().Pkg().Path(),

		"FnRegister": gengotypes.Ref("github.com/go-courier/schema/pkg/jsonschema/extractors", "Register"),
		"TSchema":    gengotypes.Ref("github.com/go-courier/schema/pkg/jsonschema", "Schema"),
		"TRefer":     gengotypes.Ref("github.com/go-courier/schema/pkg/jsonschema", "Refer"),
	}

	s := g.schemaFromType(c.Universe, named, true)

	if s == nil {
		return nil
	}

	d := gengo.NewDumper(c.Package.Pkg().Path(), g.imports)

	if err := sw.Do(`
func init() {
	{{ .FnRegister | raw }}("{{ .typePkgPath }}.{{ .typeName }}", new({{ .typeName }}))
}

func({{ .typeName }}) OpenAPISchema(ref func(t string) {{ .TRefer | raw }}) *{{ .TSchema | raw }} {
	return `+d.ValueLit(s, gengo.OnInterface(func(v interface{}) string {
		switch t := v.(type) {
		case extractors.TypeName:
			return fmt.Sprintf(`ref(%s)`, strconv.Quote(string(t)))
		}
		return d.ValueLit(v)
	}))+`
}
`, a); err != nil {
		return err
	}

	return nil
}

func (g *jsonSchemaGen) schemaFromType(u gengotypes.Universe, tpe *types.Named, def bool) (s *jsonschema.Schema) {
	docGetter := u.Package(tpe.Obj().Pkg().Path())

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

			_, doc := gengotypes.ExtractCommentTags("+", docGetter.Doc(tpe.Obj().Pos()))
			s.Description = strings.Join(doc, "\n")
		}
	}()

	ctx := extractors.WithDocGetter(context.Background(), u.Package(tpe.Obj().Pkg().Path()))

	return extractors.OpenAPISchemaGetterFromGoType(ctx, tpe).OpenAPISchema(func(t string) jsonschema.Refer {
		return extractors.TypeName(t)
	})
}
