package generators

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-courier/schema/pkg/jsonschema"
	"github.com/go-courier/schema/pkg/testutil"
	"github.com/pkg/errors"
	"go/ast"
	"io"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"strings"
)

func NewJsonSchemaGen(pkg *types.Package, args *args.GeneratorArgs) generator.Generator {
	return &jsonSchemaGen{
		DefaultGen: generator.DefaultGen{
			OptionalName: args.OutputFileBaseName + ".jsonschema",
		},
		outputPackage: pkg.Path,
		imports:       generator.NewImportTracker(),
	}
}

type jsonSchemaGen struct {
	generator.DefaultGen
	outputPackage    string
	imports          namer.ImportTracker
	typedConstValues map[string]map[*types.Type][]*types.Type
}

func (g *jsonSchemaGen) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

func (g *jsonSchemaGen) Imports(c *generator.Context) (imports []string) {
	return g.imports.ImportLines()
}

var supportedKinds = map[types.Kind]bool{
	types.Array:   true,
	types.Slice:   true,
	types.Map:     true,
	types.Alias:   true,
	types.Pointer: true,
	types.Struct:  true,
}

func (g *jsonSchemaGen) Filter(c *generator.Context, tpe *types.Type) bool {
	return ast.IsExported(tpe.Name.Name) && supportedKinds[tpe.Kind]
}

func (g *jsonSchemaGen) GenerateType(c *generator.Context, tpe *types.Type, w io.Writer) error {

	sw := generator.NewSnippetWriter(w, c, "{{", "}}")
	g.genOpenAPISchema(sw, tpe)

	return nil
}

func (g *jsonSchemaGen) schemaFromType(tpe *types.Type) (s *jsonschema.Schema) {
	defer func() {
		if s != nil {
			commentLines := make([]string, 0)

			if lines := strings.Join(tpe.CommentLines, "\n"); lines != "" {
				commentLines = append(commentLines, lines)
			}

			if lines := strings.Join(tpe.SecondClosestCommentLines, "\n"); lines != "" {
				commentLines = append(commentLines, lines)
			}

			s.Description = strings.Join(commentLines, "\n")
		}
	}()

	switch tpe.Kind {
	case types.Builtin:
		return jsonschema.NewSchema(schemaTypeAndFormatFromBasicType(tpe.Name.Name))
	case types.Alias:
		return g.schemaFromType(tpe.Underlying)
	case types.Array:
		s := jsonschema.ItemsOf(g.schemaFromType(tpe.Elem))

		spew.Dump(tpe.Underlying)

		return s
	case types.Slice:
		return jsonschema.ItemsOf(g.schemaFromType(tpe.Elem))
	case types.Map:
		keySchema := g.schemaFromType(tpe.Key)
		if keySchema != nil && len(keySchema.Type) > 0 && keySchema.Type != "string" {
			panic(errors.New("only support map[string]interface{}"))
		}
		return jsonschema.KeyValueOf(keySchema, g.schemaFromType(tpe.Elem))
	case types.Struct:

	}

	return &jsonschema.Schema{}
}

func (g *jsonSchemaGen) genOpenAPISchema(sw *generator.SnippetWriter, tpe *types.Type) {
	a := generator.Args{
		"type":                    tpe,
		"JSONSchema":              types.Ref("github.com/go-courier/schema/pkg/jsonschema", "Schema"),
		"OpenAPISchemaTypeGetter": types.Ref("github.com/go-courier/schema/pkg/jsonschema", "OpenAPISchemaTypeGetter"),
	}

	fmt.Println(tpe.Name)
	testutil.PrintJSON(g.schemaFromType(tpe))

	sw.Do(`
func (v {{ .type | raw }}) OpenAPISchema() *{{ .JSONSchema | raw }} {
	s := &{{ .JSONSchema | raw }}{}

	return s
}
`, a)
}

var basicTypeToSchemaType = map[string][2]string{
	"invalid": {"null", ""},

	"bool":    {"boolean", ""},
	"error":   {"string", "string"},
	"float32": {"number", "float"},
	"float64": {"number", "double"},

	"int":   {"integer", "int32"},
	"int8":  {"integer", "int8"},
	"int16": {"integer", "int16"},
	"int32": {"integer", "int32"},
	"int64": {"integer", "int64"},

	"rune": {"integer", "int32"},

	"uint":   {"integer", "uint32"},
	"uint8":  {"integer", "uint8"},
	"uint16": {"integer", "uint16"},
	"uint32": {"integer", "uint32"},
	"uint64": {"integer", "uint64"},

	"byte": {"integer", "uint8"},

	"string": {"string", ""},
}

func schemaTypeAndFormatFromBasicType(basicTypeName string) (typ jsonschema.Type, format string) {
	if schemaTypeAndFormat, ok := basicTypeToSchemaType[basicTypeName]; ok {
		return jsonschema.Type(schemaTypeAndFormat[0]), schemaTypeAndFormat[1]
	}
	panic(errors.Errorf("unsupported type %q", basicTypeName))
}
