package querysql

import (
	"database/sql/driver"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"

	"github.com/go-courier/gengo/pkg/gengo"
	typesx "github.com/go-courier/x/types"
)

func init() {
	gengo.Register(&QuerySqlGen{})
}

type QuerySqlGen struct {
	gengo.SnippetWriter
}

func (QuerySqlGen) Name() string {
	return "querysql"
}

func (QuerySqlGen) New() gengo.Generator {
	return &QuerySqlGen{}
}

func (g *QuerySqlGen) Init(c *gengo.Context, s gengo.GeneratorCreator) (gengo.Generator, error) {
	return s.Init(c, g, func(g gengo.Generator, sw gengo.SnippetWriter) error {
		g.(*QuerySqlGen).SnippetWriter = sw
		return nil
	})
}

func (g *QuerySqlGen) GenerateType(c *gengo.Context, named *types.Named) error {
	t := typesx.FromTType(named.Underlying())

	if t.Kind() != reflect.Struct {
		return nil
	}

	parameters := scanParameters(c, t)

	g.Do(`
func (p *[[ .typeName ]]) ToCondition(d [[ "github.com/go-courier/sqlx/v2.TableResolver" | id ]]) [[ "github.com/go-courier/sqlx/v2/builder.SqlCondition" | id ]] {
	where := [[ "github.com/go-courier/sqlx/v2/builder.EmptyCond" | id ]]()
	[[ .doEach | render ]]
	return where
}
`, gengo.Args{
		"typeName": named.Obj().Name(),
		"doEach": func(sw gengo.SnippetWriter) {
			m := map[string]bool{}

			for i := range parameters {
				p := parameters[i]

				if p.WhereModel != nil {
					model := p.WhereModel.PkgPath + "." + p.WhereModel.Name

					if _, ok := m[model]; !ok {

						g.Do(`
m[[ .name ]] := &[[ .model | id ]]{}
t[[ .name ]] := d.T(m[[ .name ]])
`, gengo.Args{
							"model": model,
							"name":  p.WhereModel.Name,
						})

						m[model] = true
					}
				}
			}

			for i := range parameters {
				p := parameters[i]

				if p.WhereModel != nil {
					if p.IsSlice {
						g.Do(`
if len(p.[[ .fieldName ]]) != 0 {
`, gengo.Args{"fieldName": p.Name})
					} else {
						g.Do(`
if ![[ "github.com/go-courier/x/reflect.IsEmptyValue" | id ]](p.[[ .fieldName ]]) {
`, gengo.Args{"fieldName": p.Name})
					}

					g.Do(`
	where = [[ "github.com/go-courier/sqlx/v2/builder.And" | id ]](
		where, 
		t[[ .whereModel.Name ]].F(m[[ .whereModel.Name ]].FieldKey[[ .whereModel.FieldName ]]()).
			In(p.[[ .fieldName ]]),
	)
}
`, gengo.Args{
						"fieldName":  p.Name,
						"whereModel": p.WhereModel,
					})
				}

			}
		},
	})

	return nil
}

var tpeDriverValuer = reflect.TypeOf((*driver.Valuer)(nil)).Elem()

func scanParameters(c *gengo.Context, t typesx.Type) []Parameter {
	parameters := make([]Parameter, 0, t.NumField())

	var walk func(t typesx.Type)

	walk = func(t typesx.Type) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			name := f.Name()

			if !ast.IsExported(name) {
				continue
			}

			p := Parameter{}
			p.Name = name

			ft := f.Type()

			if canPos, ok := f.(interface{ Pos() token.Pos }); ok {
				tags, _ := c.Package.Doc(canPos.Pos())

				if del, ok := tags["where"]; ok {
					wm, err := ParseModelField(del[0], func(localName string) string {
						for _, i := range c.Package.Pkg().Imports() {
							if i.Name() == "db" {
								return i.Path()
							}
						}
						return ""
					})

					if err != nil {
						panic(err)
					}

					p.WhereModel = wm
				}

			}

			if !ft.Implements(typesx.FromRType(tpeDriverValuer)) {
				p.IsSlice = ft.Kind() == reflect.Slice
			} else {
				if f.Anonymous() && ft.Kind() == reflect.Struct {
					//walk(ft)
				}
			}

			parameters = append(parameters, p)
		}
	}

	walk(t)

	return parameters
}

type Parameter struct {
	Name       string
	IsSlice    bool
	WhereModel *ModelField
}
