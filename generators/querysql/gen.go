package querysql

import (
	"go/types"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-courier/gengo/pkg/gengo"
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

	spew.Dump(named)

	return nil
}
