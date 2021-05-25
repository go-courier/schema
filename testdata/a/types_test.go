package a

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-courier/reflectx/typesutil"
	"github.com/go-courier/schema/pkg/validator"

	"github.com/go-courier/schema/pkg/jsonschema/extractors"
	"github.com/go-courier/schema/pkg/testutil"
)

type Struct2 Struct

func TestSchemaFrom(t *testing.T) {
	testutil.PrintJSON(extractors.SchemaFrom(&Struct{}))
	testutil.PrintJSON(extractors.SchemaFrom(&Struct2{}))
}

func Benchmark(b *testing.B) {
	v := validator.ValidatorMgrDefault.MustCompile(context.Background(), nil, typesutil.FromRType(reflect.TypeOf(&Struct2{})), nil)

	b.Run("by reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = v.Validate(Struct2{})
		}
	})

	b.Run("by generated", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = (&Struct{}).Validate()
		}
	})
}
