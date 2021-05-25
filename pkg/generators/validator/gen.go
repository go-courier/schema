package validator

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go/ast"
	"go/types"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-courier/gengo/pkg/gengo"
	"github.com/go-courier/gengo/pkg/namer"
	gentypes "github.com/go-courier/gengo/pkg/types"
	"github.com/go-courier/reflectx/typesutil"
	"github.com/go-courier/schema/pkg/validator"
)

func init() {
	gengo.Register(&validatorGen{})
}

type validatorGen struct {
	imports namer.ImportTracker
}

func (g *validatorGen) Name() string {
	return "validator"
}

func (validatorGen) New() gengo.Generator {
	return &validatorGen{
		imports: namer.NewDefaultImportTracker(),
	}
}

func (g *validatorGen) Imports(context *gengo.Context) map[string]string {
	return g.imports.Imports()
}

func (g *validatorGen) Init(c *gengo.Context, writer io.Writer) error {
	return nil
}

func (g *validatorGen) GenerateType(c *gengo.Context, named *types.Named, w io.Writer) error {
	d := gengo.NewDumper(c.Package.Pkg().Path(), g.imports)

	do := func(w io.Writer, s string) {
		_, _ = io.WriteString(w, s)
	}

	ifDo := func(w io.Writer, condition string, stmt func(w io.Writer)) {
		do(w, `if `+condition+` {
`)
		stmt(w)
		do(w, `}`)
	}

	tt := typesutil.FromTType(named.Underlying())

	_, _ = fmt.Fprintf(w, `func (v *`+named.Obj().Name()+`) Validate() error {`)

	if tt.Kind() == reflect.Struct {
		_, _ = fmt.Fprintf(w, `validate := func(v *`+named.Obj().Name()+`) error {
	errSet := `+d.Name(gentypes.Ref("github.com/go-courier/schema/pkg/validator", "NewErrorSet"))+`()
`)

		walkStruct(tt, func(sf *StructField) {
			validateField := func(w io.Writer, emptyValue string, tpe string, writeValidate func(w io.Writer, ref string)) {
				do(w, `
{ 
	var fe error
`)

				setMissingErrIfRequired := func(w io.Writer) {
					if !sf.Omitempty {
						do(w, ` else {
		fe = `+d.ValueLit(validator.MissingRequired{})+`
	}
`)
					}
				}

				doValidate := func(w io.Writer) {
					ref := tpe + "(" + strings.Repeat("*", sf.PtrLevel) + "v." + sf.Name + ")"

					if _, ok := typesutil.EncodingTextMarshalerTypeReplacer(sf.Type); ok {
						ref = d.Name(gentypes.Ref("github.com/go-courier/schema/pkg/validator", "ToMarshalledText")) + "(" + strings.Repeat("*", sf.PtrLevel) + "v." + sf.Name + ")"
						emptyValue = `""`
					}

					ifDo(w, `fv := `+ref+`; fv != `+emptyValue, func(w io.Writer) {
						do(w, "fe = ")
						writeValidate(w, "fv")
					})

					setMissingErrIfRequired(w)
				}

				if sf.PtrLevel > 0 {
					deNilCheckList := make([]string, sf.PtrLevel)

					for i := 0; i < sf.PtrLevel; i++ {
						deNilCheckList[i] = strings.Repeat("*", i) + "v." + sf.Name + `!= nil`
					}

					ifDo(w, strings.Join(deNilCheckList, " && "), doValidate)
					setMissingErrIfRequired(w)
				} else {
					doValidate(w)
				}

				do(w, `
	if fe != nil {
		errSet.AddErr(fe, `+strconv.Quote(sf.DisplayName)+`)
	}
}`)

			}

			switch sf.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				validateField(w, `0`, `int64`, func(w io.Writer, ref string) {
					_, _ = fmt.Fprintf(w, `(`+d.ValueLit(&validator.IntValidator{
					})+`).Validate(`+ref+`);`)
				})
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				validateField(w, `0`, `uint64`, func(w io.Writer, ref string) {
					_, _ = fmt.Fprintf(w, `(`+d.ValueLit(&validator.UintValidator{})+`).Validate(`+ref+`);`)
				})
			case reflect.Float64, reflect.Float32:
				validateField(w, `0`, `float64`, func(w io.Writer, ref string) {
					_, _ = fmt.Fprintf(w, `(`+d.ValueLit(&validator.FloatValidator{})+`).Validate(`+ref+`);`)
				})
			case reflect.String:
				validateField(w, `""`, `string`, func(w io.Writer, ref string) {
					_, _ = fmt.Fprintf(w, `(`+d.ValueLit(&validator.StringValidator{})+`).Validate(`+ref+`);`)
				})
			case reflect.Array:

			case reflect.Slice:
			case reflect.Struct:

			}
		}, []string{})

		_, _ = fmt.Fprintf(w, `
	
	return errSet.Err()
}`)
	}

	_, _ = fmt.Fprintf(w, `
	return validate(v)
}`)

	return nil
}

type StructField struct {
	Name        string
	DisplayName string
	Type        typesutil.Type
	Omitempty   bool
	PtrLevel    int
	Paths       []string
	Validator   validator.Validator
}

func walkStruct(s typesutil.Type, each func(sf *StructField), parents []string) {
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fieldName := field.Name()
		fieldTag := field.Tag()

		jsonTag, keepNested := fieldTag.Lookup("json")
		validateTag, needToValidator := fieldTag.Lookup("validate")

		displayName := strings.Split(jsonTag, ",")[0]
		omitempty := strings.Contains(jsonTag, "omitempty")

		if !ast.IsExported(fieldName) || displayName == "-" {
			continue
		}

		fieldType := field.Type()
		ptrLevel := 0

		for {
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
				ptrLevel++
			} else {
				break
			}
		}

		if field.Anonymous() {
			switch fieldType.Kind() {
			case reflect.Struct:
				if !keepNested {
					walkStruct(s, each, parents)
					continue
				}
			case reflect.Interface:
				continue
			}
		}

		sf := &StructField{}

		sf.Name = fieldName
		sf.Type = fieldType
		sf.PtrLevel = ptrLevel
		sf.DisplayName = displayName
		sf.Omitempty = omitempty

		if needToValidator {
			spew.Dump(validateTag)
		}

		each(sf)
	}
}
