package generators

import (
	"io"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"strconv"
	"strings"
)

func NewEnumGen(pkg *types.Package, args *args.GeneratorArgs) generator.Generator {
	return &enumGen{
		DefaultGen: generator.DefaultGen{
			OptionalName: args.OutputFileBaseName + ".enum",
		},
		outputPackage: pkg.Path,
		imports:       generator.NewImportTracker(),
	}
}

type enumGen struct {
	generator.DefaultGen
	outputPackage    string
	imports          namer.ImportTracker
	typedConstValues map[string]map[*types.Type][]*types.Type
}

func (g *enumGen) Init(c *generator.Context, w io.Writer) error {
	g.typedConstValues = map[string]map[*types.Type][]*types.Type{}

	for i := range c.Universe {
		p := c.Universe[i]

		for k := range p.Constants {
			if g.typedConstValues[p.Path] == nil {
				g.typedConstValues[p.Path] = map[*types.Type][]*types.Type{}
			}
			constant := p.Constants[k]
			g.typedConstValues[p.Path][constant.Underlying] = append(g.typedConstValues[p.Path][constant.Underlying], constant)
		}
	}

	return g.DefaultGen.Init(c, w)
}

func (g *enumGen) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

func (g *enumGen) Imports(c *generator.Context) (imports []string) {
	return g.imports.ImportLines()
}

func (g *enumGen) GenerateType(c *generator.Context, tpe *types.Type, w io.Writer) error {
	pkg := tpe.Name.Package

	if typedConstValues, ok := g.typedConstValues[pkg]; ok {
		if constValues := typedConstValues[tpe]; ok {
			if len(constValues) > 0 {
				sw := generator.NewSnippetWriter(w, c, "{{", "}}")

				_, e := strconv.ParseInt(*constValues[0].ConstValue, 10, 64)

				// for int constants
				if e == nil {
					constValues = g.genIntStringerMethods(sw, tpe, constValues)
				}

				g.genOpenAPISchemaEnum(sw, tpe, constValues)
			}
		}
	}

	return nil
}

type option struct {
	Int         int64
	Label       string
	QuotedLabel string
	Str         string
	QuotedStr   string
	Type        *types.Type
}

func (g *enumGen) genOpenAPISchemaEnum(sw *generator.SnippetWriter, tpe *types.Type, constValues []*types.Type) {
	a := generator.Args{
		"type":        tpe,
		"constValues": constValues,
	}

	sw.Do(`
func ({{ .type | raw }}) OpenAPISchemaEnum() []interface{} {
	return []interface{}{
		{{ range .constValues }} {{ . | raw }}, 
		{{ end }} }
}
`, a)
}

func (g *enumGen) genIntStringerMethods(sw *generator.SnippetWriter, tpe *types.Type, constValues []*types.Type) (validConstValues []*types.Type) {
	options := make([]*option, 0)

	var constUnknown *types.Type

	for i := range constValues {
		cv := constValues[i]
		n := cv.Name.Name
		if n[0] == '_' {
			continue
		}

		if strings.HasSuffix(n, "_UNKNOWN") {
			constUnknown = cv
		}

		var values = strings.SplitN(n, "__", 2)
		if len(values) == 2 {
			o := &option{}

			o.Str = values[1]
			o.Label = strings.Join(cv.TrailingCommentLines, "")
			o.Int, _ = strconv.ParseInt(*cv.ConstValue, 10, 64)
			o.Type = cv

			if o.Label == "" {
				o.Label = o.Str
			}

			o.QuotedStr = strconv.Quote(o.Str)
			o.QuotedLabel = strconv.Quote(o.Label)

			options = append(options, o)
			validConstValues = append(validConstValues, o.Type)
		}
	}

	a := generator.Args{
		"type":         tpe,
		"constUnknown": constUnknown,
		"options":      options,

		"NewError":       types.Ref("github.com/pkg/errors", "New"),
		"SqlDriverValue": types.Ref("database/sql/driver", "Value"),

		"IntStringerEnum":     types.Ref("github.com/go-courier/schema/pkg/enumeration", "IntStringerEnum"),
		"ScanIntEnumStringer": types.Ref("github.com/go-courier/schema/pkg/enumeration", "ScanIntEnumStringer"),
		"DriverValueOffset":   types.Ref("github.com/go-courier/schema/pkg/enumeration", "DriverValueOffset"),
	}

	sw.Do(`
var Invalid{{ .type | raw }} = {{ .NewError | raw }}("invalid {{ .type | raw }}")

func Parse{{ .type | raw }}FromString(s string) ({{ .type | raw }}, error) {
	switch s {
	{{ range .options }} case {{ .QuotedStr }}:
		return {{ .Type | raw }}, nil 
	{{ end }} }
	return {{ .constUnknown | raw }}, Invalid{{ .type | raw }}
}

func Parse{{ .type | raw }}FromLabelString(s string) ({{ .type | raw }}, error) {
	switch s {
	{{ range .options }} case {{ .QuotedLabel }}:
		return {{ .Type | raw }}, nil
	{{ end }} }
	return {{ .constUnknown | raw }}, Invalid{{ .type | raw }}
}

func ({{ .type | raw }}) TypeName() string {
	return "{{ .type.Name.Package }}.{{ .type.Name.Name }}"
}

func (v {{ .type | raw }}) String() string {
	switch v {
	{{ range .options }} case {{ .Type.Name.Name }}:
		return {{ .QuotedStr }} 
	{{ end }} }
	return "UNKNOWN"
}


func (v {{ .type | raw }}) Label() string {
	switch v {
	{{ range .options }} case {{ .Type | raw }}:
		return {{ .QuotedLabel }} 
	{{ end }} }
	return "UNKNOWN"
}

func (v {{ .type | raw }}) Int() int {
	return int(v)
}

func ({{ .type | raw }}) ConstValues() []{{ .IntStringerEnum | raw }} {
	return []{{ .IntStringerEnum | raw }}{
		{{ range .options }}{{ .Type.Name.Name }}, 
		{{ end }} }
}

func (v {{ .type | raw }}) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, Invalid{{ .type | raw }}
	}
	return []byte(str), nil
}

func (v *{{ .type | raw }}) UnmarshalText(data []byte) (err error) {
	*v, err = Parse{{ .type | raw }}FromString(string(bytes.ToUpper(data)))
	return
}


func (v {{ .type | raw }}) Value() ({{ .SqlDriverValue | raw }}, error) {
	offset := 0
	if o, ok := (interface{})(v).({{ .DriverValueOffset | raw }}); ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}

func (v *{{ .type | raw }}) Scan(src interface{}) error {
	offset := 0
	if o, ok := (interface{})(v).({{ .DriverValueOffset | raw }}); ok {
		offset = o.Offset()
	}

	i, err := {{ .ScanIntEnumStringer | raw }}(src, offset)
	if err != nil {
		return err
	}
	*v = {{ .type | raw }}(i)
	return nil
}
`, a)

	return
}
