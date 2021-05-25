package openapi

import (
	"context"
	"go/types"

	"github.com/go-courier/gengo/pkg/gengo"
	"github.com/go-courier/httptransport/transformers"
	typesx "github.com/go-courier/x/types"
)

func (g *OperatorGen) generateFromRequest(c *gengo.Context, named *types.Named) error {
	ctx := context.Background()
	t := typesx.FromTType(named)

	groupedParameters := map[string][]*transformers.RequestParameter{}
	hasParameters := false

	if err := transformers.EachRequestParameter(ctx, t, func(p *transformers.RequestParameter) {
		if p.In != "" {
			groupedParameters[p.In] = append(groupedParameters[p.In], p)
			hasParameters = true
		}
	}); err != nil {
		return err
	}

	g.Do(`
func(r *[[ .typeName ]]) FromRequestInfo(ri [[ "github.com/go-courier/httptransport/httpx.RequestInfo" | id ]]) error {
[[ if .hasParameters ]] errSet := [[ "github.com/go-courier/schema/pkg/validator.NewErrorSet" | id ]]()
	[[ .unmarshalFromRequest | render ]]
	if errSet.Len() == 0 {
		return nil
	}
	return errSet
[[ else ]] return nil [[ end ]]
}
`, gengo.Args{
		"hasParameters": hasParameters,
		"typeName":      named.Obj().Name(),
		"unmarshalFromRequest": func(s gengo.SnippetWriter) {
			for _, loc := range []string{"path", "query", "header", "cookie", "body"} {
				if parameters, ok := groupedParameters[loc]; ok {
					for i := range parameters {
						p := parameters[i]

						if p.In == "body" {
							s.Do(`
body := ri.Body()
if err := ([[ .transformer ]]).DecodeFrom(ri.Context(), body, &r.[[ .fieldName ]], [[ "net/textproto.MIMEHeader" | id ]](ri.Header())); err != nil {
	errSet.AddErr(err, [[ "github.com/go-courier/schema/pkg/validator.Location" | id ]]([[ .param.In | quote ]]))
}
body.Close()
`, gengo.Args{
								"param":             p,
								"fieldName":         p.Field.Name(),
								"transformer":       s.Dumper().ValueLit(p.Transformer),
								"transformerOption": s.Dumper().ValueLit(&p.TransformerOption.CommonTransformOption),
							})

							continue
						}

						if p.TransformerOption.CommonTransformOption.Explode {
							s.Do(`
if values := ri.Values([[ .param.In | quote ]], [[ .param.Name | quote ]]); len(values) > 0 {
	[[ if eq .fieldType "[]string" ]] r.[[ .fieldName ]] = values [[ else ]]
	n := len(values)
	
	r.[[ .fieldName ]] = make([[ .fieldType ]], n)

	for i := range r.[[ .fieldName ]] {	
		if err := ([[ .transformer ]]).DecodeFrom(ri.Context(), [[ "github.com/go-courier/httptransport/transformers.NewStringReader" | id ]](values[i]), &r.[[ .fieldName ]][i]); err != nil {
			errSet.AddErr(err, [[ "github.com/go-courier/schema/pkg/validator.Location" | id ]]([[ .param.In | quote ]]), [[ .param.Name | quote ]], i)
		}
	}
	[[ end ]]
}
`, gengo.Args{
								"param":       p,
								"fieldType":   g.Dumper().TypeLit(p.Field.Type()),
								"fieldName":   p.Field.Name(),
								"transformer": s.Dumper().ValueLit(p.Transformer),
							})
						} else {
							s.Do(`
if values := ri.Values([[ .param.In | quote ]], [[ .param.Name | quote ]]); len(values) > 0 {
	[[ if eq .fieldType "string" ]] r.[[ .fieldName ]] = values[0] [[ else ]]
	if err := ([[ .transformer ]]).DecodeFrom(ri.Context(), [[ "github.com/go-courier/httptransport/transformers.NewStringReader" | id ]](values[0]), &r.[[ .fieldName ]]); err != nil {
		errSet.AddErr(err, [[ "github.com/go-courier/schema/pkg/validator.Location" | id ]]([[ .param.In | quote ]]), [[ .param.Name | quote ]])
	}
	[[ end ]]
}
`, gengo.Args{
								"param":       p,
								"fieldType":   p.Field.Type().String(),
								"fieldName":   p.Field.Name(),
								"transformer": s.Dumper().ValueLit(p.Transformer),
							})
						}
					}
				}
			}
		},
	})

	return nil
}
