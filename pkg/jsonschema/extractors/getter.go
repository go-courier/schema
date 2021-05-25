package extractors

import (
	"context"
	"go/types"
	"reflect"

	"github.com/go-courier/reflectx/typesutil"
	"github.com/go-courier/schema/pkg/jsonschema"
)

func OpenAPISchemaGetterFromReflectType(ctx context.Context, t reflect.Type) OpenAPISchemaGetter {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return &valueOpenAPISchemaGetter{
		ctx: ctx,
		t:   typesutil.FromRType(t),
	}
}

func OpenAPISchemaGetterFromGoType(ctx context.Context, types types.Type, def bool) OpenAPISchemaGetter {
	return &valueOpenAPISchemaGetter{
		t:   typesutil.FromTType(types),
		def: def,
		ctx: ctx,
	}
}

type valueOpenAPISchemaGetter struct {
	ctx context.Context
	t   typesutil.Type
	def bool
}

func (v *valueOpenAPISchemaGetter) OpenAPISchema(ref func(t string) jsonschema.Refer) *jsonschema.Schema {
	return SchemaFromType(v.ctx, v.t, v.def)
}

type opt struct {
	Ref func(t string) jsonschema.Refer
}

type SchemaFromOptFn = func(o *opt)

func SchemaFromOptRef(ref func(t string) jsonschema.Refer) SchemaFromOptFn {
	return func(o *opt) {
		o.Ref = ref
	}
}

func SchemaFrom(v interface{}, optionFns ...SchemaFromOptFn) (s *jsonschema.Schema) {
	o := opt{}

	for i := range optionFns {
		optionFns[i](&o)
	}

	if o.Ref == nil {
		o.Ref = func(t string) jsonschema.Refer {
			return TypeName(t)
		}
	}

	defer func() {
		if g, ok := v.(OpenAPISchemaTypeGetter); ok {
			s.Type = g.OpenAPISchemaType()
			s.Format = ""
		}

		if g, ok := v.(OpenAPISchemaFormatGetter); ok {
			s.Format = g.OpenAPISchemaFormat()
		}
	}()

	if g, ok := v.(OpenAPISchemaGetter); ok {
		return g.OpenAPISchema(o.Ref)
	}

	return OpenAPISchemaGetterFromReflectType(context.Background(), reflect.TypeOf(v)).OpenAPISchema(o.Ref)
}
