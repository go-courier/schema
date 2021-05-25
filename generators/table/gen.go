package Table

import (
	"context"
	"fmt"
	"github.com/go-courier/gengo/pkg/gengo"
	"github.com/go-courier/sqlx/v2/builder"
	typesx "github.com/go-courier/x/types"
	"go/types"
	"strings"
)

func init() {
	gengo.Register(&TableGen{})
}

type TableGen struct {
	gengo.SnippetWriter
}

func (TableGen) Name() string {
	return "table"
}

func (TableGen) New() gengo.Generator {
	return &TableGen{}
}

func (g *TableGen) Init(c *gengo.Context, s gengo.GeneratorCreator) (gengo.Generator, error) {
	return s.Init(c, g, func(g gengo.Generator, sw gengo.SnippetWriter) error {
		g.(*TableGen).SnippetWriter = sw
		return nil
	})
}

func toDefaultTableName(name string) string {
	return gengo.LowerSnakeCase("t_" + name)
}

func (g *TableGen) GenerateType(c *gengo.Context, named *types.Named) error {
	t := g.scanTable(c, named)

	g.generateIndexInterfaces(t, named)
	g.generateDescriptions(t, named)

	return nil
}

func (g *TableGen) generateDescriptions(t *builder.Table, named *types.Named) {
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

func (g *TableGen) generateIndexInterfaces(t *builder.Table, named *types.Named) {
	primary := make([]string, 0)
	indexes := builder.Indexes{}
	uniqueIndexes := builder.Indexes{}

	t.Keys.Range(func(key *builder.Key, idx int) {
		if key.IsPrimary() {
			primary = key.Columns.FieldNames()
		} else {
			n := key.Name
			if key.Method != "" {
				n = n + "/" + key.Method
			}
			if key.IsUnique {
				uniqueIndexes[n] = key.Columns.FieldNames()
			} else {
				indexes[n] = key.Columns.FieldNames()
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

func (g *TableGen) scanTable(c *gengo.Context, named *types.Named) *builder.Table {
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
			defs := defSplit(indexes[i])

			switch strings.TrimSpace(defs[0]) {
			case "primary":
				if len(defs) < 2 {
					panic(fmt.Errorf("primary at lease 1 Field"))
				}
				cols, err := t.Fields(defs[1:]...)
				if err != nil {
					panic(fmt.Errorf("%s, please check primary def", err))
				}
				t.AddKey(builder.PrimaryKey(cols))

			case "index":
				if len(defs) < 3 {
					panic(fmt.Errorf("index at lease 1 StructField"))
				}
				indexNameAndMethod := defs[1]
				fieldNames := defs[2:]

				indexName, method := builder.ResolveIndexNameAndMethod(indexNameAndMethod)
				cols, err := t.Fields(fieldNames...)
				if err != nil {
					panic(fmt.Errorf("%s, please check index def", err))
				}
				t.AddKey(builder.Index(indexName, cols).Using(method))

			case "unique_index":
				if len(defs) < 3 {
					panic(fmt.Errorf("unique indexes at lease 1 StructField"))
				}
				indexNameAndMethod := defs[1]
				fieldNames := defs[2:]

				indexName, method := builder.ResolveIndexNameAndMethod(indexNameAndMethod)
				cols, err := t.Fields(fieldNames...)
				if err != nil {
					panic(fmt.Errorf("%s, please check unique_index def", err))
				}
				t.AddKey(builder.UniqueIndex(indexName, cols).Using(method))

			}
		}
	}

	return t
}

func defSplit(def string) (defs []string) {
	vs := strings.Split(def, " ")
	for _, s := range vs {
		if s != "" {
			defs = append(defs, s)
		}
	}
	return
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
