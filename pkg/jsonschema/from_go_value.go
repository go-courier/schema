package jsonschema

import (
	"encoding"
	"errors"
	"fmt"
	"github.com/go-courier/codegen"
	"go/ast"
	"mime/multipart"
	"reflect"
	"strings"
)

const (
	XEnumLabels   = `x-enum-labels`
	XGoVendorType = `x-go-vendor-type`
)

type OpenAPISchemaTypeGetter interface {
	OpenAPISchemaType() []string
}

type OpenAPISchemaFormatGetter interface {
	OpenAPISchemaFormat() string
}

type OpenAPISchemaEnumGetter interface {
	OpenAPISchemaEnum() []interface{}
}

func FromGoValue(v interface{}) *Schema {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	return (&SchemaPicker{
		NamedTags:   []string{"json", "name"},
		Definitions: map[string]*Schema{},
	}).Pick(rv)
}

type SchemaModifier func(s *Schema, rv reflect.Value)

func schemaModifierNullable(s *Schema, rv reflect.Value) {
	if s != nil {
		s.Nullable = true
	}
}

func schemaModifierMayEnum(s *Schema, rv reflect.Value) {
	if enumGetter, ok := rv.Interface().(OpenAPISchemaEnumGetter); ok {
		enum := enumGetter.OpenAPISchemaEnum()

		enumOptions := make([]string, len(enum))

		for i := range enum {
			if l, ok := enum[i].(interface{ Label() string }); ok {
				enumOptions[i] = l.Label()
			} else {
				enumOptions[i] = fmt.Sprintf("%v", enum[i])
			}
		}

		s.Enum = enum

		s.AddExtension(XEnumLabels, enumOptions)
	}
}

func schemaModifierMayFormat(s *Schema, rv reflect.Value) {
	if enumGetter, ok := rv.Interface().(OpenAPISchemaFormatGetter); ok {
		s.Format = enumGetter.OpenAPISchemaFormat()
	}
}

type SchemaPicker struct {
	NamedTags   []string
	Definitions map[string]*Schema
}

func (c *SchemaPicker) Pick(rv reflect.Value) (s *Schema) {
	return c.pick(rv, schemaModifierMayEnum, schemaModifierMayFormat)
}

type OpenAPISchemaDescription string

func (d OpenAPISchemaDescription) OpenAPISchemaDescription(fieldPath ...string) string {
	return string(d)
}

func (c *SchemaPicker) pick(rv reflect.Value, schemaModifies ...SchemaModifier) (s *Schema) {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv = reflect.New(rv.Type().Elem())
		}
		return c.pick(rv.Elem(), append(schemaModifies, schemaModifierNullable)...)
	}

	t := rv.Type()
	goVendorType := ""

	if pkgPath := t.PkgPath(); pkgPath != "" {
		goVendorType = fullTypeName(t)

		id := codegen.UpperCamelCase(goVendorType)

		if _, ok := c.Definitions[id]; ok {
			return RefSchema("#/definitions/" + id)
		}

		c.Definitions[id] = s
	}

	defer func() {
		if s != nil {
			for i := range schemaModifies {
				schemaModifies[i](s, rv)
			}

			if goVendorType != "" {
				s.AddExtension(XGoVendorType, goVendorType)
			}
		}
	}()

	if _, ok := rv.Interface().(encoding.TextMarshaler); ok {
		return String()
	}

	if _, ok := rv.Interface().(multipart.FileHeader); ok {
		// todo ugly hard code
		return Binary()
	}

	switch rv.Kind() {
	case reflect.Bool:
		return Boolean()
	case reflect.String:
		return String()
	case reflect.Float64:
		return Double()
	case reflect.Float32:
		return Float()
	case reflect.Int:
		return NewSchema(TypeInteger)
	case reflect.Int8:
		return NewSchema(TypeInteger, "int8")
	case reflect.Int16:
		return NewSchema(TypeInteger, "int16")
	case reflect.Int32:
		return NewSchema(TypeInteger, "int32")
	case reflect.Int64:
		return NewSchema(TypeInteger, "int64")
	case reflect.Uint:
		return NewSchema(TypeInteger, "uint")
	case reflect.Uint8:
		return NewSchema(TypeInteger, "uint8")
	case reflect.Uint16:
		return NewSchema(TypeInteger, "uint16")
	case reflect.Uint32:
		return NewSchema(TypeInteger, "uint32")
	case reflect.Uint64:
		return NewSchema(TypeInteger, "uint64")
	case reflect.Interface:
		return &Schema{}
	case reflect.Map:
		keySchema := c.Pick(reflect.New(t.Key()).Elem())
		if keySchema != nil && len(keySchema.Type) > 0 && keySchema.Type != "string" {
			panic(errors.New("only support map[string]interface{}"))
		}
		return KeyValueOf(keySchema, c.Pick(reflect.New(t.Elem()).Elem()))
	case reflect.Slice:
		return ItemsOf(c.Pick(reflect.New(t.Elem()).Elem()))
	case reflect.Array:
		length := uint64(t.Len())
		s := ItemsOf(c.Pick(reflect.New(t.Elem()).Elem()))
		s.MaxItems = &length
		s.MinItems = &length
		return s
	case reflect.Struct:
		structSchema := ObjectOf(nil)

		schemas := make([]*Schema, 0)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			if !ast.IsExported(field.Name) {
				continue
			}

			fieldValue := rv.Field(i)

			tags := field.Tag

			tagValueForName := ""

			for _, namedTag := range c.NamedTags {
				if tagValueForName == "" {
					tagValueForName = tags.Get(namedTag)
				}
			}

			name, flags := tagValueAndFlagsByTagString(tagValueForName)
			if name == "-" {
				continue
			}

			if name == "" && field.Anonymous {
				if field.Type.String() == "bytes.Buffer" {
					// todo fix
					structSchema = Binary()
					break
				}
				s := c.Pick(fieldValue)
				if s != nil {
					schemas = append(schemas, s)
				}
				continue
			}

			if name == "" {
				name = field.Name
			}

			required := true

			if hasOmitempty, ok := flags["omitempty"]; ok {
				required = !hasOmitempty
			}

			s := c.Pick(fieldValue)

			structSchema.SetProperty(name, s, required)
		}

		if len(schemas) > 0 {
			return AllOf(append(schemas, structSchema)...)
		}

		return structSchema
	}

	return nil
}

func tagValueAndFlagsByTagString(tagString string) (string, map[string]bool) {
	valueAndFlags := strings.Split(tagString, ",")
	v := valueAndFlags[0]
	tagFlags := map[string]bool{}
	if len(valueAndFlags) > 1 {
		for _, flag := range valueAndFlags[1:] {
			tagFlags[flag] = true
		}
	}
	return v, tagFlags
}

func fullTypeName(typeName reflect.Type) string {
	pkgPath := typeName.PkgPath()
	if pkgPath != "" {
		return pkgPath + "." + typeName.Name()
	}
	return typeName.Name()
}
