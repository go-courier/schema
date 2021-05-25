package openapi

import (
	"context"
	"go/ast"
	"go/constant"
	"go/types"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-courier/courier"
	"github.com/go-courier/gengo/pkg/gengo"
	gengotypes "github.com/go-courier/gengo/pkg/types"
	"github.com/go-courier/httptransport/httpx"
	"github.com/go-courier/httptransport/transformers"
	"github.com/go-courier/reflectx/typesutil"
	jsonschemagenerator "github.com/go-courier/schema/pkg/generators/jsonschema"
	"github.com/go-courier/schema/pkg/jsonschema"
	"github.com/go-courier/schema/pkg/jsonschema/extractors"
	"github.com/go-courier/schema/pkg/openapi"
	"github.com/go-courier/statuserror"
	"github.com/pkg/errors"
)

func init() {
	gengo.Register(&OperatorGen{})
}

type OperatorGen struct {
	gengo.SnippetWriter

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
		return g.generateOpenAPIOperation(c, named)
	}

	return g.jsonschemaGen.GenerateType(c, named)
}

func (g *OperatorGen) generateOpenAPIOperation(c *gengo.Context, named *types.Named) error {
	operation := &openapi.Operation{}

	if tags, ok := c.Tags["gengo:"+g.Name()+":tag"]; ok {
		operation.Tags = tags
	}

	operation.OperationId = named.Obj().Name()

	docGetter := c.Universe.Package(named.Obj().Pkg().Path())
	ctx := extractors.WithDocGetter(context.Background(), docGetter)

	tags, lines := docGetter.Doc(named.Obj().Pos())

	if _, ok := tags["deprecated"]; ok {
		operation.Deprecated = true
	}

	if len(lines) > 1 {
		operation.Summary = lines[0]
		if len(lines) > 2 {
			operation.Description = strings.Join(lines[1:], "\n")
		}
	}

	path := g.scanParametersOrRequestBody(c, ctx, operation, typesutil.FromTType(named.Underlying().(*types.Struct)), func(t string) jsonschema.Refer {
		return extractors.TypeName(t)
	})

	g.scanOutputReturns(c, operation, named)

	operationLit, err := g.jsonschemaGen.ValueLitAndGenSideTypes(c, g.Dumper(), operation)
	if err != nil {
		return err
	}

	g.Do(`
func([[ .typeName ]]) New() [[ .TOperator | raw ]] {
	return &[[ .typeName ]]{}
}
[[ if ne .path "" ]]
func([[ .typeName ]]) Path() string {
	return [[ .path | quote ]]
} [[ end ]]

func([[ .typeName ]]) OpenAPIOperation(ref func(t string) [[ .TRefer | raw ]]) *[[ .TOperation | raw ]] {
	return [[ .operation ]]
}
`, gengo.Args{
		"path":       path,
		"TOperator":  gengotypes.Ref("github.com/go-courier/courier", "Operator"),
		"TOperation": gengotypes.Ref("github.com/go-courier/schema/pkg/openapi", "Operation"),
		"TRefer":     gengotypes.Ref("github.com/go-courier/schema/pkg/jsonschema", "Refer"),

		"typeName":  named.Obj().Name(),
		"operation": operationLit,
	})

	return nil
}

func (g *OperatorGen) scanOutputReturns(c *gengo.Context, o *openapi.Operation, named *types.Named) {
	method, ok := typesutil.FromTType(types.NewPointer(named)).MethodByName("Output")
	if ok {
		results, n := c.Package.ResultsOf(method.(*typesutil.TMethod).Func)
		if n == 2 {
			httpMethod := http.MethodGet

			if m, ok := g.firstValueOfFunc(c, named, "Method"); ok {
				httpMethod = m.(string)
			}

			g.setResponseDataFromResults(c, o, results[0], httpMethod)
			g.setResponseStatusErrorsFromResults(c, o, results[1])
		}
	}
}

func (g *OperatorGen) setResponseDataFromResults(c *gengo.Context, o *openapi.Operation, typeAndValues []gengotypes.TypeAndValue, httpMethod string) {
	statusCode := http.StatusNoContent

	response := &openapi.Response{}
	response.Description = http.StatusText(statusCode)

	for _, resp := range typeAndValues {
		if resp.Type != nil && !isNil(resp.Type) {
			if "UpdateByID" == o.OperationId {
				spew.Dump(o.OperationId, resp.Type)
			}
			statusCode, response = g.resolveResponse(c, resp.Type, resp.Expr, httpMethod)
		}
	}

	o.Responses.AddResponse(statusCode, response)
}

func (g *OperatorGen) setResponseStatusErrorsFromResults(c *gengo.Context, o *openapi.Operation, typeAndValues []gengotypes.TypeAndValue) {
	statusErrors := map[int][]*statuserror.StatusErr{}

	for _, err := range typeAndValues {
		if err.Type != nil && err.Type.String() != types.Typ[types.UntypedNil].String() {
			if isStatusErr(err.Type) {
				switch x := err.Expr.(type) {
				case *ast.CallExpr:
					identList := getIdentChainOfCallFunc(x)

					if len(identList) > 0 {
						callIdent := identList[len(identList)-1]

						// pick status errors from github.com/go-courier/statuserror.Wrap
						if callIdent.Name == "Wrap" {
							code := 0
							key := ""
							msg := ""
							desc := make([]string, 0)

							for i, arg := range x.Args[1:] {
								p := c.Universe.LocateInPackage(x.Pos())
								if p == nil {
									continue
								}

								tv, err := p.Eval(arg)
								if err != nil {
									continue
								}

								switch i {
								case 0: // code
									code, _ = strconv.Atoi(tv.Value.String())
								case 1: // key
									key, _ = strconv.Unquote(tv.Value.String())
								case 2: // msg
									msg, _ = strconv.Unquote(tv.Value.String())
								default:
									d, _ := strconv.Unquote(tv.Value.String())
									desc = append(desc, d)
								}
							}

							if code > 0 {
								if msg == "" {
									msg = key
								}

								statusErrors[code] = append(statusErrors[code], statuserror.Wrap(errors.New(""), code, key, append([]string{msg}, desc...)...))
							}
						}
					}
				}
			}
		}
	}

	for statusCode := range statusErrors {
		response := openapi.NewResponse("")

		statusErrs := make([]string, len(statusErrors[statusCode]))

		for i, s := range statusErrors[statusCode] {
			statusErrs[i] = s.Summary()
		}

		sort.Strings(statusErrs)

		response.AddExtension(XStatusErrs, statusErrs)

		tpeStatusErr := c.Universe.Package(pkgImportPathStatusError).Type("StatusErr")

		named := tpeStatusErr.Type().(*types.Named)

		s := g.jsonschemaGen.SchemaFromType(context.Background(), tpeStatusErr.Pos(), named, false)

		response.AddContent(httpx.MIME_JSON, openapi.NewMediaTypeWithSchema(s))

		o.Responses.AddResponse(statusCode, response)
	}
}

func (g *OperatorGen) resolveResponse(c *gengo.Context, tpe types.Type, expr ast.Expr, httpMethod string) (statusCode int, response *openapi.Response) {
	defer func() {
		if statusCode == 0 {
			switch httpMethod {
			case http.MethodPost:
				statusCode = http.StatusCreated
			default:
				statusCode = http.StatusOK
			}
		}
	}()

	response = &openapi.Response{}
	response.Description = "-"
	contentType := ""

	if isHttpxResponse(tpe) {
		scanResponseWrapper := func(expr ast.Expr) {
			firstCallExpr := true

			ast.Inspect(expr, func(node ast.Node) bool {
				switch callExpr := node.(type) {
				case *ast.CallExpr:
					if firstCallExpr {
						firstCallExpr = false
						if p := c.Universe.LocateInPackage(node.Pos()); p != nil {
							v, _ := p.Eval(callExpr.Args[0])
							tpe = v.Type
						}
					}
					switch e := callExpr.Fun.(type) {
					case *ast.SelectorExpr:
						switch e.Sel.Name {
						case "WithSchema":
							if p := c.Universe.LocateInPackage(node.Pos()); p != nil {
								v, _ := p.Eval(callExpr.Args[0])
								tpe = v.Type
							}
						case "WithStatusCode":
							if p := c.Universe.LocateInPackage(node.Pos()); p != nil {
								v, _ := p.Eval(callExpr.Args[0])
								if code, ok := valueOf(v.Value).(int64); ok {
									statusCode = int(code)
								}
							}
							return false
						case "WithContentType":
							if p := c.Universe.LocateInPackage(node.Pos()); p != nil {
								v, _ := p.Eval(callExpr.Args[0])
								if ct, ok := valueOf(v.Value).(string); ok {
									contentType = ct
								}
							}
							return false
						}
					}
				}
				return true
			})
		}

		if ident, ok := expr.(*ast.Ident); ok && ident.Obj != nil {
			if stmt, ok := ident.Obj.Decl.(*ast.AssignStmt); ok {
				for _, e := range stmt.Rhs {
					scanResponseWrapper(e)
				}
			}
		} else {
			scanResponseWrapper(expr)
		}
	}

	if pointer, ok := tpe.(*types.Pointer); ok {
		tpe = pointer.Elem()
	}

	if named, ok := tpe.(*types.Named); ok {
		if v, ok := g.firstValueOfFunc(c, named, "ContentType"); ok {
			if s, ok := v.(string); ok {
				contentType = s
			}
			if contentType == "" {
				contentType = "*"
			}
		}
		if v, ok := g.firstValueOfFunc(c, named, "StatusCode"); ok {
			if i, ok := v.(int64); ok {
				statusCode = int(i)
			}
		}
	}

	// skip redirect
	if statusCode > http.StatusMultipleChoices && statusCode < http.StatusBadRequest {
		response.Description = http.StatusText(statusCode)
		return
	}

	if isNil(tpe) {
		statusCode = http.StatusNoContent
		return
	}

	if contentType == "" {
		contentType = httpx.MIME_JSON
	}

	s := g.jsonschemaGen.SchemaFromType(context.Background(), expr.Pos(), tpe, false)

	response.AddContent(contentType, openapi.NewMediaTypeWithSchema(s))

	return
}

func (g *OperatorGen) scanParametersOrRequestBody(c *gengo.Context, ctx context.Context, op *openapi.Operation, s typesutil.Type, ref func(t string) jsonschema.Refer) (path string) {
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		structTag := field.Tag()

		pathTag, hasPath := structTag.Lookup("path")
		if hasPath {
			path = pathTag
		}

		nameTag, hasNameTag := structTag.Lookup("name")

		if field.Anonymous() && !hasNameTag {
			if ft := field.Type(); ft.Kind() == reflect.Struct {
				g.scanParametersOrRequestBody(c, ctx, op, ft, ref)
				continue
			}
		}

		fieldDisplayName := field.Name()
		omitempty := false

		if nameTag != "" {
			if n := strings.Split(nameTag, ",")[0]; n != "" {
				fieldDisplayName = n
			}
			omitempty = strings.Contains(nameTag, "omitempty")
		}

		location := structTag.Get("in")

		if location == "" {
			panic(errors.Errorf("missing tag `in` for %s of %s", field.Name(), op.OperationId))
		}

		schema := extractors.PropSchemaFromStructField(
			ctx,
			field,
			!omitempty,
		)

		transformer, err := transformers.TransformerMgrDefault.NewTransformer(context.Background(), field.Type(), transformers.TransformerOption{
			MIME: structTag.Get("mime"),
		})

		if err != nil {
			panic(err)
		}

		switch location {
		case "body":
			reqBody := openapi.NewRequestBody("", true)
			reqBody.AddContent(transformer.Names()[0], openapi.NewMediaTypeWithSchema(schema))
			op.SetRequestBody(reqBody)
		case "query":
			op.AddParameter(openapi.QueryParameter(fieldDisplayName, schema, !omitempty))
		case "cookie":
			op.AddParameter(openapi.CookieParameter(fieldDisplayName, schema, !omitempty))
		case "header":
			op.AddParameter(openapi.HeaderParameter(fieldDisplayName, schema, !omitempty))
		case "path":
			op.AddParameter(openapi.PathParameter(fieldDisplayName, schema))
		}
	}

	return
}

func (g *OperatorGen) firstValueOfFunc(c *gengo.Context, named *types.Named, name string) (interface{}, bool) {
	method, ok := typesutil.FromTType(types.NewPointer(named)).MethodByName(name)
	if ok {
		fn := method.(*typesutil.TMethod).Func
		results, n := c.Universe.Package(fn.Pkg().Path()).ResultsOf(fn)
		if n == 1 {
			for _, r := range results[0] {
				if v := valueOf(r.Value); v != nil {
					return v, true
				}
			}
			return nil, true
		}
	}
	return nil, false
}

func valueOf(v constant.Value) interface{} {
	if v == nil {
		return nil
	}

	switch v.Kind() {
	case constant.Float:
		v, _ := strconv.ParseFloat(v.String(), 10)
		return v
	case constant.Bool:
		v, _ := strconv.ParseBool(v.String())
		return v
	case constant.String:
		v, _ := strconv.Unquote(v.String())
		return v
	case constant.Int:
		v, _ := strconv.ParseInt(v.String(), 10, 64)
		return v
	}

	return nil
}
