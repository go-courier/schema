package openapi

import (
	"go/ast"
	"go/types"
	"reflect"

	"github.com/go-courier/courier"
	"github.com/go-courier/gengo/pkg/gengo"
	jsonschemagenerator "github.com/go-courier/schema/generators/jsonschema"
	"github.com/go-courier/schema/generators/validator"
	typesutil "github.com/go-courier/x/types"
)

func init() {
	gengo.Register(&OperatorGen{})
}

type OperatorGen struct {
	gengo.SnippetWriter

	validatorGen  *validator.ValidatorGen
	jsonschemaGen *jsonschemagenerator.JSONSchemaGen
}

func (g *OperatorGen) Name() string {
	return "operator"
}

func (OperatorGen) New() gengo.Generator {
	return &OperatorGen{}
}

func (g *OperatorGen) Init(c *gengo.Context, s gengo.GeneratorCreator) (gengo.Generator, error) {
	return s.Init(c, g, func(g gengo.Generator, sw gengo.SnippetWriter) error {
		g.(*OperatorGen).SnippetWriter = sw

		if jsonschemaGen, err := (&jsonschemagenerator.JSONSchemaGen{}).Init(c, s); err != nil {
			return err
		} else {
			g.(*OperatorGen).jsonschemaGen = jsonschemaGen.(*jsonschemagenerator.JSONSchemaGen)
		}

		if validatorGen, err := (&validator.ValidatorGen{}).Init(c, s); err != nil {
			return err
		} else {
			g.(*OperatorGen).validatorGen = validatorGen.(*validator.ValidatorGen)
		}

		return nil
	})
}

var typOperator = reflect.TypeOf((*courier.Operator)(nil)).Elem()

func isCourierOperator(tpe typesutil.Type, lookup func(importPath string) *types.Package) bool {
	switch tpe.(type) {
	case *typesutil.RType:
		return tpe.Implements(typesutil.FromRType(typOperator))
	case *typesutil.TType:
		pkg := lookup(typOperator.PkgPath())
		if pkg == nil {
			return false
		}
		t := pkg.Scope().Lookup(typOperator.Name())
		if t == nil {
			return false
		}
		return types.Implements(tpe.(*typesutil.TType).Type, t.Type().Underlying().(*types.Interface))
	}
	return false
}

func (g *OperatorGen) GenerateType(c *gengo.Context, named *types.Named) error {
	// skip private
	if !ast.IsExported(named.Obj().Name()) {
		return nil
	}

	ptrTyp := typesutil.FromTType(types.NewPointer(named))

	if isCourierOperator(ptrTyp, func(importPath string) *types.Package { return c.Universe.Package(importPath).Pkg() }) {
		if err := g.validatorGen.GenerateType(c, named); err != nil {
			return err
		}
		if err := g.generateFromRequest(c, named); err != nil {
			return err
		}
		return g.generateOpenAPIOperation(c, named)
	}

	return g.jsonschemaGen.GenerateType(c, named)
}
