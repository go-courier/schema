package gocourier

import (
	"context"
	"net/http"
	"regexp"
	"sort"

	"github.com/go-courier/courier"
	"github.com/go-courier/gengo/pkg/gengo"
	"github.com/go-courier/httptransport"
	"github.com/go-courier/httptransport/httpx"
	"github.com/go-courier/schema/pkg/jsonschema"
	"github.com/go-courier/schema/pkg/jsonschema/extractors"
	"github.com/go-courier/schema/pkg/openapi"
)

type OpenAPIPostProcess = func(o *openapi.OpenAPI)

func OpenAPIFor(rootRouter *courier.Router, postProcesses ...OpenAPIPostProcess) *courier.Router {
	return courier.NewRouter(&OpenAPISpec{root: rootRouter, postProcesses: postProcesses})
}

type OpenAPISpec struct {
	httpx.MethodGet

	root          *courier.Router
	postProcesses []OpenAPIPostProcess
}

func (req *OpenAPISpec) New() courier.Operator {
	return &OpenAPISpec{
		root:          req.root,
		postProcesses: req.postProcesses,
	}
}

func (req *OpenAPISpec) Output(c context.Context) (interface{}, error) {
	return resolveOpenAPI(req.root, req.postProcesses...), nil
}

func resolveOpenAPI(rootRouter *courier.Router, postProcesses ...OpenAPIPostProcess) *openapi.OpenAPI {
	o := openapi.NewOpenAPI()

	defer func() {
		for i := range postProcesses {
			postProcesses[i](o)
		}
	}()

	components := componentSchemas{}

	o.Schemas = components

	routes := rootRouter.Routes()

	routeMetas := make([]*httptransport.HttpRouteMeta, len(routes))

	for i := range routes {
		routeMetas[i] = httptransport.NewHttpRouteMeta(routes[i])
	}

	sort.Slice(routeMetas, func(i, j int) bool {
		return routeMetas[i].Key() < routeMetas[j].Key()
	})

	for _, httpRoute := range routeMetas {
		operation := &openapi.Operation{}

		parameters := make([]*openapi.Parameter, 0)

		for _, of := range httpRoute.OperatorFactoryWithRouteMetas {
			o := &openapi.Operation{}

			og, ok := of.Operator.(openapi.OperationGetter)
			if ok {
				o = og.OpenAPIOperation(components.Ref)
			} else {
				o.OperationId = of.ID

				r := &openapi.Response{}
				r.Description = "-"
				o.Responses.AddResponse(http.StatusOK, r)
			}

			parameters = append(parameters, o.Parameters...)

			if of.IsLast {
				operation = o

				if len(o.Tags) == 0 {
					o.Tags = []string{"internal"}
				}

				sort.Slice(parameters, func(i, j int) bool {
					return parameters[i].Name < parameters[j].Name
				})

				operation.Parameters = parameters
			}
		}

		o.AddOperation(openapi.HttpMethod(httpRoute.Method()), reformatPath(httpRoute.Path()), operation)
	}

	return o
}

type componentSchemas map[string]*openapi.Schema

func (d componentSchemas) Ref(importPathAndName string) jsonschema.Refer {
	n := gengo.UpperCamelCase(importPathAndName)

	if g, ok := extractors.Get(importPathAndName); ok {
		if _, ok := d[n]; !ok {
			// register first to avoid loop issue
			d[n] = &jsonschema.Schema{}

			d[n] = extractors.SchemaFrom(g, extractors.SchemaFromOptRef(d.Ref))
		}
	}

	return jsonschema.RefSchema("#/components/schemas/" + n)
}

var re = regexp.MustCompile("/:([^/]+)")

func reformatPath(path string) string {
	return re.ReplaceAllString(path, "/{$1}")
}
