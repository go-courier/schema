package Table

import (
	"context"
	"go/types"
	"strings"

	"github.com/go-courier/gengo/pkg/gengo"
	"github.com/go-courier/sqlx/v2/builder"
	typesx "github.com/go-courier/x/types"
)

func init() {
	gengo.Register(&TableModelGen{})
}

type TableModelGen struct {
	gengo.SnippetWriter
}

func (TableModelGen) Name() string {
	return "tablemodel"
}

func (TableModelGen) New() gengo.Generator {
	return &TableModelGen{}
}

func (g *TableModelGen) Init(c *gengo.Context, s gengo.GeneratorCreator) (gengo.Generator, error) {
	return s.Init(c, g, func(g gengo.Generator, sw gengo.SnippetWriter) error {
		g.(*TableModelGen).SnippetWriter = sw
		return nil
	})
}

func toDefaultTableName(name string) string {
	return gengo.LowerSnakeCase("t_" + name)
}

func (g *TableModelGen) GenerateType(c *gengo.Context, named *types.Named) error {
	t := g.scanTable(c, named)

	g.generateConvertors(t, named)
	g.generateIndexInterfaces(t, named)
	g.generateFieldKeys(t, named)

	return nil
}

func (g *TableModelGen) generateConvertors(t *builder.Table, named *types.Named) {
	g.Do(`
func(v *[[ .typeName ]]) ColumnReceivers() map[string]interface{} {
	return map[string]interface{}{
		[[ .columnReceivers | render ]]
	}
}
`, gengo.Args{
		"typeName": named.Obj().Name(),
		"columnReceivers": func(sw gengo.SnippetWriter) {
			t.Columns.Range(func(col *builder.Column, idx int) {
				sw.Do(`[[ .name | quote ]]: &v.[[ .fieldName ]],
`, gengo.Args{
					"name":      col.Name,
					"fieldName": col.FieldName,
				})
			})
		},
	})
}

func (g *TableModelGen) generateFieldKeys(t *builder.Table, named *types.Named) {
	t.Columns.Range(func(col *builder.Column, idx int) {
		if col.DeprecatedActions == nil {
			g.Do(`
func([[ .typeName ]]) FieldKey[[ .fieldName ]]() string {
	return [[ .fieldName | quote ]]
}
`, gengo.Args{
				"typeName":  named.Obj().Name(),
				"fieldName": col.FieldName,
			})
		}

	})
}

func (g *TableModelGen) generateDescriptions(t *builder.Table, named *types.Named) {
	colComments := map[string]string{}
	colDescriptions := map[string][]string{}
	colRelations := map[string][]string{}

	t.Columns.Range(func(col *builder.Column, idx int) {
		if col.Comment != "" {
			colComments[col.FieldName] = col.Comment
		}
		if len(col.Description) > 0 {
			colDescriptions[col.FieldName] = col.Description
		}
		if len(col.Relation) > 0 {
			colRelations[col.FieldName] = col.Relation
		}
	})

	g.Do(`
[[ if .hasTableDescription ]] func([[ .typeName ]]) TableDescription() []string {
	return [[ .tableDescription ]]
} [[ end ]]

[[ if .hasComments ]] func([[ .typeName ]]) Comments() map[string]string {
	return [[ .comments ]]
} [[ end ]]

[[ if .hasColDescriptions ]] func([[ .typeName ]]) ColDescriptions() map[string][]string {
	return [[ .colDescriptions ]]
} [[ end ]]


[[ if .hasColRelations ]] func([[ .typeName ]]) ColRelations() map[string][]string {
	return [[ .colRelations ]]
} [[ end ]]
`, gengo.Args{
		"typeName": named.Obj().Name(),

		"hasTableDescription": len(t.Description) > 0,
		"tableDescription":    g.Dumper().ValueLit(t.Description),

		"hasComments": len(colComments) > 0,
		"comments":    g.Dumper().ValueLit(colComments),

		"hasColDescriptions": len(colDescriptions) > 0,
		"colDescriptions":    g.Dumper().ValueLit(colDescriptions),

		"hasColRelations": len(colRelations) > 0,
		"colRelations":    g.Dumper().ValueLit(colRelations),
	})
}

func (g *TableModelGen) generateIndexInterfaces(t *builder.Table, named *types.Named) {
	primary := make([]string, 0)
	indexes := builder.Indexes{}
	uniqueIndexes := builder.Indexes{}

	t.Keys.Range(func(key *builder.Key, idx int) {
		if key.IsPrimary() {
			primary = key.Def.ToDefs()
		} else {
			n := key.Name
			if key.Method != "" {
				n = n + "/" + key.Method
			}
			if key.IsUnique {
				uniqueIndexes[n] = key.Def.ToDefs()
			} else {
				indexes[n] = key.Def.ToDefs()
			}
		}
	})

	g.Do(`
func ([[ .typeName ]]) TableName() string {
	return [[ .tableName | quote ]]
}

[[ if .hasPrimary ]] func([[ .typeName ]]) Primary() []string {
	return [[ .primary ]]
} [[ end ]]

[[ if .hasIndexes ]] func([[ .typeName ]]) Indexes() [[ "github.com/go-courier/sqlx/v2/builder.Indexes" | id ]] {
	return [[ .indexes ]]
} [[ end ]]

[[ if .hasUniqueIndexes ]] func([[ .typeName ]]) UniqueIndexes() [[ "github.com/go-courier/sqlx/v2/builder.Indexes" | id ]] {
	return [[ .uniqueIndexes ]]
} [[ end ]]
`, gengo.Args{
		"typeName":  named.Obj().Name(),
		"tableName": t.Name,

		"hasPrimary":       len(primary) > 0,
		"primary":          g.Dumper().ValueLit(primary),
		"hasUniqueIndexes": len(uniqueIndexes) > 0,
		"uniqueIndexes":    g.Dumper().ValueLit(uniqueIndexes),
		"hasIndexes":       len(indexes) > 0,
		"indexes":          g.Dumper().ValueLit(indexes),
	})
}

func (g *TableModelGen) scanTable(c *gengo.Context, named *types.Named) *builder.Table {
	tags, doc := c.Universe.Package(named.Obj().Pkg().Path()).Doc(named.Obj().Pos())

	tableName := toDefaultTableName(named.Obj().Name())
	if tn, ok := tags["gengo:table:name"]; ok {
		if n := tn[0]; len(n) > 0 {
			tableName = n
		}
	}

	t := builder.T(tableName)

	t.Description = doc

	builder.EachStructField(context.Background(), typesx.FromTType(named), func(p *builder.StructField) bool {
		col := builder.Col(p.Name).Field(p.FieldName)

		if tsf, ok := p.Field.(*typesx.TStructField); ok {
			var tags map[string][]string
			var doc []string

			if pkgPath := p.Field.PkgPath(); pkgPath != "" {
				tags, doc = c.Universe.Package(pkgPath).Doc(tsf.Pos())
			} else {
				tags, doc = c.Universe.Package(named.Obj().Pkg().Path()).Doc(tsf.Pos())
			}

			col.Comment, col.Description = commentAndDesc(doc)

			if values, ok := tags["rel"]; ok {
				rel := strings.Split(values[0], ".")
				if len(rel) >= 2 {
					col.Relation = rel
				}
			}
		}

		t.AddCol(col)
		return true
	})

	if indexes, ok := tags["def"]; ok {
		for i := range indexes {
			def := builder.ParseIndexDefine(indexes[i])

			switch def.Kind {
			case "primary":
				t.AddKey(&builder.Key{
					Name:     def.Kind,
					IsUnique: true,
					Def:      def.IndexDef,
				})
			case "index":
				t.AddKey(&builder.Key{
					Name:   def.Name,
					Method: def.Method,
					Def:    def.IndexDef,
				})
			case "unique_index":
				t.AddKey(&builder.Key{
					Name:     def.Name,
					Method:   def.Method,
					IsUnique: true,
					Def:      def.IndexDef,
				})
			}
		}
	}

	return t
}

func commentAndDesc(docs []string) (comment string, desc []string) {
	for _, s := range docs {
		if comment == "" && s == "" {
			continue
		}
		if comment == "" {
			comment = s
		} else {
			desc = append(desc, s)
		}
	}
	return
}
