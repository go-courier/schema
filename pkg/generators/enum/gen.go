package enum

import (
	"fmt"
	"go/types"
	"strconv"

	"github.com/go-courier/gengo/pkg/gengo"
	gengotypes "github.com/go-courier/gengo/pkg/types"
)

func init() {
	gengo.Register(&EnumGen{})
}

type EnumGen struct {
	gengo.SnippetWriter
	enumTypes EnumTypes
}

func (EnumGen) Name() string {
	return "enum"
}

func (EnumGen) New() gengo.Generator {
	return &EnumGen{
		enumTypes: EnumTypes{},
	}
}

func (g *EnumGen) Init(c *gengo.Context, s gengo.GeneratorCreator) (gengo.Generator, error) {
	return s.Init(c, g, func(g gengo.Generator, sw gengo.SnippetWriter) error {
		g.(*EnumGen).SnippetWriter = sw
		g.(*EnumGen).enumTypes.Walk(c, c.Package.Pkg().Path())
		return nil
	})
}

func (g *EnumGen) GenerateType(c *gengo.Context, named *types.Named) error {
	if enum, ok := g.enumTypes.ResolveEnumType(named); ok {
		if enum.IsIntStringer() {
			g.genIntStringerMethods(named, enum)
		}
	}
	return nil
}

func (g *EnumGen) genIntStringerMethods(tpe types.Type, enum *EnumType) {
	options := make([]struct {
		Name        string
		QuotedValue string
		QuotedLabel string
	}, len(enum.Constants))

	tpeObj := tpe.(*types.Named).Obj()

	for i := range enum.Constants {
		options[i].Name = enum.Constants[i].Name()
		options[i].QuotedValue = strconv.Quote(fmt.Sprintf("%v", enum.Value(enum.Constants[i])))
		options[i].QuotedLabel = strconv.Quote(enum.Label(enum.Constants[i]))
	}

	a := gengo.Args{
		"typeName":    tpeObj.Name(),
		"typePkgPath": tpeObj.Pkg().Path(),

		"constUnknown": enum.ConstUnknown,
		"options":      options,

		"ToUpper":        gengotypes.Ref("bytes", "ToUpper"),
		"NewError":       gengotypes.Ref("github.com/pkg/errors", "New"),
		"SqlDriverValue": gengotypes.Ref("database/sql/driver", "Value"),

		"IntStringerEnum":     gengotypes.Ref("github.com/go-courier/schema/pkg/enumeration", "IntStringerEnum"),
		"ScanIntEnumStringer": gengotypes.Ref("github.com/go-courier/schema/pkg/enumeration", "ScanIntEnumStringer"),
		"DriverValueOffset":   gengotypes.Ref("github.com/go-courier/schema/pkg/enumeration", "DriverValueOffset"),
	}

	g.Do(`
var Invalid[[ .typeName ]] = [[ .NewError | raw ]]("invalid [[ .typeName ]]")

func Parse[[ .typeName ]]FromString(s string) ([[ .typeName ]], error) {
	switch s {
	[[ range .options ]] case [[ .QuotedValue ]]:
		return [[ .Name ]], nil 
	[[ end ]] }
	return [[ .constUnknown | raw ]], Invalid[[ .typeName ]]
}

func Parse[[ .typeName ]]FromLabelString(s string) ([[ .typeName ]], error) {
	switch s {
	[[ range .options ]] case [[ .QuotedLabel ]]:
		return [[ .Name ]], nil
	[[ end ]] }
	return [[ .constUnknown | raw ]], Invalid[[ .typeName ]]
}

func ([[ .typeName ]]) TypeName() string {
	return "[[ .typePkgPath ]].[[ .typeName ]]"
}

func (v [[ .typeName ]]) String() string {
	switch v {
	[[ range .options ]] case [[ .Name ]]:
		return [[ .QuotedValue ]] 
	[[ end ]] }
	return "UNKNOWN"
}


func (v [[ .typeName ]]) Label() string {
	switch v {
	[[ range .options ]] case [[ .Name ]]:
		return [[ .QuotedLabel ]] 
	[[ end ]] }
	return "UNKNOWN"
}

func (v [[ .typeName ]]) Int() int {
	return int(v)
}

func ([[ .typeName ]]) ConstValues() [][[ .IntStringerEnum | raw ]] {
	return [][[ .IntStringerEnum | raw ]]{
		[[ range .options ]][[ .Name ]], 
		[[ end ]] }
}

func (v [[ .typeName ]]) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, Invalid[[ .typeName ]]
	}
	return []byte(str), nil
}

func (v *[[ .typeName ]]) UnmarshalText(data []byte) (err error) {
	*v, err = Parse[[ .typeName ]]FromString(string([[ .ToUpper | raw ]](data)))
	return
}


func (v [[ .typeName ]]) Value() ([[ .SqlDriverValue | raw ]], error) {
	offset := 0
	if o, ok := (interface{})(v).([[ .DriverValueOffset | raw ]]); ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}

func (v *[[ .typeName ]]) Scan(src interface{}) error {
	offset := 0
	if o, ok := (interface{})(v).([[ .DriverValueOffset | raw ]]); ok {
		offset = o.Offset()
	}

	i, err := [[ .ScanIntEnumStringer | raw ]](src, offset)
	if err != nil {
		return err
	}
	*v = [[ .typeName ]](i)
	return nil
}
`, a)
}
