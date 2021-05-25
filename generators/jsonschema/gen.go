package jsonschema

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"github.com/go-courier/gengo/pkg/gengo"
	gengotypes "github.com/go-courier/gengo/pkg/types"
	"github.com/go-courier/schema/pkg/jsonschema/extractors"
)

func init() {
	gengo.Register(&JSONSchemaGen{})
}

type JSONSchemaGen struct {
	gengo.SnippetWriter
	Scanner
	processed map[string]bool
}

func (g *JSONSchemaGen) Name() string {
	return "jsonschema"
}

func (JSONSchemaGen) New() gengo.Generator {
	return &JSONSchemaGen{
		processed: map[string]bool{},
	}
}

func (g *JSONSchemaGen) Init(c *gengo.Context, s gengo.GeneratorCreator) (gengo.Generator, error) {
	return s.Init(c, g, func(g gengo.Generator, sw gengo.SnippetWriter) error {
		g.(*JSONSchemaGen).SnippetWriter = sw
		g.(*JSONSchemaGen).Scanner.Init(c)
		return nil
	})
}

func (g *JSONSchemaGen) GenerateType(c *gengo.Context, named *types.Named) error {
	// skip private
	if !ast.IsExported(named.Obj().Name()) {
		return nil
	}

	return g.GenType(c, &Type{
		PkgPath: named.Obj().Pkg().Path(),
		Name:    named.Obj().Name(),
		Pos:     named.Obj().Pos(),
		T:       named,
	})
}

func (g *JSONSchemaGen) GenType(c *gengo.Context, t *Type) error {
	id := t.ID()

	if _, processed := g.processed[id]; processed {
		return nil
	}

	g.processed[id] = true

	args := gengo.Args{
		"id":       id,
		"typeName": t.Name,
	}

	docGetter := c.Universe.Package(t.RealPkgPath())
	ctx := extractors.WithDocGetter(context.Background(), docGetter)

	s := g.SchemaFromType(ctx, t.RealPos(), t.RealType(), true)

	if s == nil {
		return nil
	}

	if t.Shadowed != nil {
		g.Do(`
type [[ .typeName ]] struct {
	[[ .type | id ]]
}
`, gengo.Args{
			"typeName": t.Name,
			"type":     gengotypes.Ref(t.Shadowed.PkgPath, t.Shadowed.Name),
		})
	}

	schemaValueLit, err := g.ValueLitAndGenSideTypes(c, g.Dumper(), s)
	if err != nil {
		return err
	}

	g.Do(`
func init() {
	[[ "github.com/go-courier/schema/pkg/jsonschema/extractors.Register" | id ]]([[ .id | quote ]], new([[ .typeName ]]))
}

func([[ .typeName ]]) OpenAPISchema(ref func(t string) [[ "github.com/go-courier/schema/pkg/jsonschema.Refer" | id ]]) *[[ "github.com/go-courier/schema/pkg/jsonschema.Schema" | id ]] {
	return [[ .openAPISchema ]]
}

`, args, gengo.Args{
		"openAPISchema": schemaValueLit,
	})

	return nil
}

func (g *JSONSchemaGen) ValueLitAndGenSideTypes(c *gengo.Context, d *gengo.Dumper, s interface{}) (string, error) {
	vendorTypes := map[string]*Type{}

	schema := d.ValueLit(s, gengo.OnInterface(func(v interface{}) string {
		switch t := v.(type) {
		case extractors.TypeName:
			if m := c.Package.Module(); m != nil {
				// to gen vendor types
				if !strings.HasPrefix(string(t), m.Path) {
					importPath, exported := importPathAndExported(string(t))
					tpe := c.Universe.Package(importPath).Type(exported)
					if tpe != nil {
						vendorTypes[string(t)] = &Type{
							PkgPath: c.Package.Pkg().Path(),
							Name:    "shadowed" + exported,
							Shadowed: &Type{
								PkgPath: importPath,
								Name:    exported,
								Pos:     tpe.Pos(),
								T:       tpe.Type(),
							},
						}
					}
				}
			}
			return fmt.Sprintf(`ref(%s)`, strconv.Quote(string(t)))
		}
		return d.ValueLit(v)
	}))

	if n := len(vendorTypes); n > 0 {
		importPaths := make([]string, n)
		i := 0

		for importPath := range vendorTypes {
			importPaths[i] = importPath
			i++
		}

		sort.Strings(importPaths)

		for _, k := range importPaths {
			if err := g.GenType(c, vendorTypes[k]); err != nil {
				return "", err
			}
		}
	}

	return schema, nil
}

func importPathAndExported(id string) (string, string) {
	parts := strings.Split(id, ".")
	if len(parts) == 1 {
		return "", parts[0]
	}
	n := len(parts)
	return strings.Join(parts[0:n-1], "."), parts[n-1]
}

type Type struct {
	PkgPath  string
	Name     string
	Pos      token.Pos
	T        types.Type
	Shadowed *Type
}

func (t *Type) ID() string {
	if t.Shadowed != nil {
		return t.Shadowed.ID()
	}
	return t.PkgPath + "." + t.Name
}

func (t *Type) RealPkgPath() string {
	if t.Shadowed != nil {
		return t.Shadowed.PkgPath
	}
	return t.PkgPath
}

func (t *Type) RealType() types.Type {
	if t.Shadowed != nil {
		return t.Shadowed.T
	}
	return t.T
}

func (t *Type) RealPos() token.Pos {
	if t.Shadowed != nil {
		return t.Shadowed.Pos
	}
	return t.Pos
}
