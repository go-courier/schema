package validator

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-courier/gengo/pkg/gengo"
	gengotypes "github.com/go-courier/gengo/pkg/types"
	"github.com/go-courier/reflectx/typesutil"
	"github.com/go-courier/schema/pkg/validator"
)

func init() {
	gengo.Register(&ValidatorGen{})
}

type ValidatorGen struct {
	gengo.SnippetWriter

	generated map[string]bool
}

func (g *ValidatorGen) Name() string {
	return "validator"
}

func (g *ValidatorGen) New() gengo.Generator {
	return &ValidatorGen{}
}

func (g *ValidatorGen) Init(c *gengo.Context, s gengo.GeneratorCreator) (gengo.Generator, error) {
	return s.Init(c, g, func(g gengo.Generator, sw gengo.SnippetWriter) error {
		g.(*ValidatorGen).SnippetWriter = sw
		g.(*ValidatorGen).generated = map[string]bool{}
		return nil
	})
}

func validateFnName(paths ...string) string {
	return strings.Join(paths, "_")
}

func (g *ValidatorGen) GenerateType(c *gengo.Context, named *types.Named) error {
	tt := typesutil.FromTType(named)

	if err := g.scanAndGenValidateFn(c, tt, &validator.ValidatorLoader{Validator: &validator.StructValidator{}}, nil); err != nil {
		return err
	}

	g.Do(`
func (v *[[ .typeName ]]) Validate() error {
	return [[ .validateFnName ]](v)
}
`, gengo.Args{
		"typeName":       named.Obj().Name(),
		"validateFnName": validateFnName("validate", named.Obj().Name()),
	})

	return nil
}

type ValidateFn struct {
	Name string
	Type typesutil.Type
}

const (
	NamePtrElem = "Ptr"
	NameKey     = "Key"
	NameElem    = "Elem"
	NameField   = "Field"
)

func (g *ValidatorGen) scanAndGenValidateFn(c *gengo.Context, tt typesutil.Type, vt *validator.ValidatorLoader, paths []string) (err error) {
	var fnName = ""

	if len(paths) == 0 {
		fnName = validateFnName("validate", tt.Name())
	} else {
		fnName = validateFnName(paths...)
	}

	if _, ok := g.generated[fnName]; ok {
		return nil
	}

	g.generated[fnName] = true

	g.genValidateFn(c, fnName, tt, vt)(g)

	switch tt.Kind() {
	case reflect.Ptr:
		if e := g.scanAndGenValidateFn(c, tt.Elem(), vt, append(paths, NamePtrElem)); e != nil {
			return e
		}
	case reflect.Struct:
		if pkgPath := tt.PkgPath(); pkgPath != "" {
			vFnName := validateFnName("validate", tt.Name())

			if vFnName != fnName {
				// skip generated named struct
				if _, ok := g.generated[vFnName]; ok {
					return nil
				} else {
					// make sure related validate fn gen
					if err := g.scanAndGenValidateFn(c, tt, vt, nil); err != nil {
						return err
					}
				}
			}
		}

		walkStruct(tt, func(sf *StructField) {
			if sf.Validator == nil {
				return
			}
			if e := g.scanAndGenValidateFn(c, sf.Type, sf.Validator, []string{fnName, NameField, sf.Name}); e != nil {
				err = e
			}
		}, paths)
		return
	case reflect.Map:
		mvt := vt.Validator.(*validator.MapValidator)

		if keyValidator := mvt.KeyValidator; keyValidator != nil {
			if e := g.scanAndGenValidateFn(c, tt.Key(), keyValidator.(*validator.ValidatorLoader), []string{fnName, NameKey}); e != nil {
				return e
			}
		}

		if elemValidator := mvt.ElemValidator; elemValidator != nil {
			if e := g.scanAndGenValidateFn(c, tt.Elem(), elemValidator.(*validator.ValidatorLoader), []string{fnName, NameElem}); e != nil {
				return e
			}
		}

	case reflect.Slice:
		mvt := vt.Validator.(*validator.SliceValidator)

		if elemValidator := mvt.ElemValidator; elemValidator != nil {
			if e := g.scanAndGenValidateFn(c, tt.Elem(), elemValidator.(*validator.ValidatorLoader), []string{fnName, NameElem}); e != nil {
				return e
			}
		}

	}

	return
}

func (g *ValidatorGen) genValidateFn(c *gengo.Context, fnName string, tt typesutil.Type, vt *validator.ValidatorLoader) func(sw gengo.SnippetWriter) {
	d := g.Dumper()

	return gengo.Snippet(`
func [[ .fnName ]](v *[[ .typeName ]]) error {
	[[ .doValidate | render ]]  
}
`, gengo.Args{
		"fnName":   fnName,
		"typeName": d.TypeLit(tt),
		"doValidate": func(sw gengo.SnippetWriter) {
			if vt.Validator == nil {
				sw.Do("return nil")
				return
			}

			if _, ok := typesutil.EncodingTextMarshalerTypeReplacer(tt); ok {
				sw.Do(`
vv, _ := v.MarshalText()
if vv != nil {
	return ([[ .vt ]]).Validate(string(vv));	
}
[[ if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
`, gengo.Args{
					"vt":         d.ValueLit(vt.Validator),
					"required":   !vt.Optional,
					"missingErr": d.ValueLit(&validator.MissingRequired{}),
				})
				return
			}

			if pkgPath := tt.PkgPath(); pkgPath != "" && tt.Kind() == reflect.Struct {
				vFnName := validateFnName("validate", tt.Name())
				// use the generated named struct
				if vFnName != fnName {
					sw.Do("return [[ .validateFnName ]](v)", gengo.Args{
						"validateFnName": vFnName,
					})
					return
				}
			}

			switch tt.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				sw.Do(`
vv := *v
if vv == 0 {
	[[ if .shouldSetDefault ]] vv = [[ .defaultValue ]]; *v = vv; [[ else if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
}
return ([[ .vt ]]).Validate(int64(vv));	
`, gengo.Args{
					"shouldSetDefault": len(vt.DefaultValue) != 0,
					"defaultValue":     string(vt.DefaultValue),
					"vt":               d.ValueLit(vt.Validator),
					"required":         !vt.Optional,
					"missingErr":       d.ValueLit(&validator.MissingRequired{}),
				})
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				sw.Do(`
vv := *v
if vv == 0 {
	[[ if .shouldSetDefault ]] vv = [[ .defaultValue ]]; *v = vv; [[ else if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
}
return ([[ .vt ]]).Validate(uint64(vv));	
`, gengo.Args{
					"shouldSetDefault": len(vt.DefaultValue) != 0,
					"defaultValue":     string(vt.DefaultValue),
					"vt":               d.ValueLit(vt.Validator),
					"required":         !vt.Optional,
					"missingErr":       d.ValueLit(&validator.MissingRequired{}),
				})
			case reflect.Float64, reflect.Float32:
				sw.Do(`
vv := *v
if vv == 0 {
	[[ if .shouldSetDefault ]] vv = [[ .defaultValue ]]; *v = vv; [[ else if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
}
return ([[ .vt ]]).Validate(float64(vv));
`, gengo.Args{
					"shouldSetDefault": len(vt.DefaultValue) != 0,
					"defaultValue":     string(vt.DefaultValue),
					"vt":               d.ValueLit(vt.Validator),
					"required":         !vt.Optional,
					"missingErr":       d.ValueLit(&validator.MissingRequired{}),
				})
			case reflect.String:
				sw.Do(`
vv := *v
if vv == "" {
	[[ if .shouldSetDefault ]] vv = [[ .defaultValue ]]; *v = vv; [[ else if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
}
return ([[ .vt ]]).Validate(string(vv));
`, gengo.Args{
					"shouldSetDefault": len(vt.DefaultValue) != 0,
					"defaultValue":     strconv.Quote(string(vt.DefaultValue)),
					"vt":               d.ValueLit(vt.Validator),
					"required":         !vt.Optional,
					"missingErr":       d.ValueLit(&validator.MissingRequired{}),
				})
			case reflect.Ptr:
				sw.Do(`
vv := *v
if vv == nil {
	[[ if .shouldSetDefault ]] var vvv [[ .typeName ]]; vv = &vvv; *v = vv; [[ else if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
}
return [[ .validateElemFnName ]](vv);	
`, gengo.Args{
					"shouldSetDefault":   len(vt.DefaultValue) != 0,
					"required":           !vt.Optional,
					"validateElemFnName": validateFnName(fnName, NamePtrElem),
					"missingErr":         d.ValueLit(&validator.MissingRequired{}),
					"typeName":           d.TypeLit(tt.Elem()),
				})
			case reflect.Map:
				mapValidator := vt.Validator.(*validator.MapValidator)

				sw.Do(`
m := *v;
`)

				if mapValidator.MinProperties > 0 {
					sw.Do(`
if n := len(m); n < [[ .minProperties ]] {
	return [[ .err ]]
}
`, gengo.Args{
						"minProperties": mapValidator.MinProperties,
						"err": d.ValueLit(&validator.OutOfRangeError{
							Current: "CURRENT",
							Minimum: mapValidator.MinProperties,
							Target:  validator.TargetSliceLength,
						}, gengo.OnInterface(func(v interface{}) string {
							if v == "CURRENT" {
								return "n"
							}
							return fmt.Sprintf("%v", v)
						})),
					})
				}

				if mapValidator.MaxProperties != nil {
					sw.Do(`
if n := len(m); n > [[ .maxItems ]] {
	return [[ .err ]]
}
`, gengo.Args{
						"maxItems": *mapValidator.MaxProperties,
						"err": d.ValueLit(&validator.OutOfRangeError{
							Current: "CURRENT",
							Maximum: *mapValidator.MaxProperties,
							Target:  validator.TargetSliceLength,
						}, gengo.OnInterface(func(v interface{}) string {
							if v == "CURRENT" {
								return "n"
							}
							return fmt.Sprintf("%v", v)
						})),
					})
				}

				args := gengo.Args{
					"FnNewErrSet":        gengotypes.Ref("github.com/go-courier/schema/pkg/validator", "NewErrorSet"),
					"validateElemFnName": validateFnName(fnName, NameElem),
					"shouldValidateKey":  mapValidator.KeyValidator != nil,
				}

				if mapValidator.KeyValidator != nil {
					args["validateKeyFnName"] = validateFnName(fnName, NameKey)
				}

				sw.Do(`
errSet := [[ .FnNewErrSet | raw ]]()
for k := range m {
	[[ if .shouldValidateKey ]] if e := [[ .validateKeyFnName ]](&k); e != nil {
			errSet.AddErr(e, k)
	} [[ end ]]
	value := m[k]
	if e := [[ .validateElemFnName ]](&value); e != nil {
		errSet.AddErr(e, k)
	}
}
if errSet.Len() > 0 {
	return errSet
}
return nil
`, args)

			case reflect.Array, reflect.Slice:
				sliceValidator := vt.Validator.(*validator.SliceValidator)

				sw.Do(`
list := *v;
n := len(list);
`)

				if sliceValidator.MinItems > 0 {
					sw.Do(`
if n < [[ .minItems ]] {
	return [[ .err ]]
}
`, gengo.Args{
						"minItems": sliceValidator.MinItems,
						"err": d.ValueLit(&validator.OutOfRangeError{
							Current: "CURRENT",
							Minimum: sliceValidator.MinItems,
							Target:  validator.TargetSliceLength,
						}, gengo.OnInterface(func(v interface{}) string {
							if v == "CURRENT" {
								return "n"
							}
							return fmt.Sprintf("%v", v)
						})),
					})
				}

				if sliceValidator.MaxItems != nil {
					sw.Do(`
if n > [[ .maxItems ]] {
	return [[ .err ]]
}
`, gengo.Args{
						"maxItems": *sliceValidator.MaxItems,
						"err": d.ValueLit(&validator.OutOfRangeError{
							Current: "CURRENT",
							Maximum: *sliceValidator.MaxItems,
							Target:  validator.TargetSliceLength,
						}, gengo.OnInterface(func(v interface{}) string {
							if v == "CURRENT" {
								return "n"
							}
							return fmt.Sprintf("%v", v)
						})),
					})
				}

				sw.Do(`
errSet := [[ .FnNewErrSet | raw ]]()
for i := 0; i < n; i++ {
	errSet.AddErr([[ .validateElemFnName ]](&list[i]), i)
}
if errSet.Len() > 0 {
	return errSet
}
return nil
`, gengo.Args{
					"FnNewErrSet":        gengotypes.Ref("github.com/go-courier/schema/pkg/validator", "NewErrorSet"),
					"validateElemFnName": validateFnName(fnName, NameElem),
				})

			case reflect.Struct:
				sw.Do(`
errSet := [[ .FnNewErrSet | raw ]]()
[[ .validateFields | render ]] if errSet.Len() > 0 {
	return errSet
}
return nil
`, gengo.Args{
					"FnNewErrSet": gengotypes.Ref("github.com/go-courier/schema/pkg/validator", "NewErrorSet"),
					"validateFields": func(sw gengo.SnippetWriter) {
						walkStruct(tt, func(sf *StructField) {
							if sf.Validator == nil {
								return
							}

							sw.Do(`
errSet.AddErr([[ .validateFnName ]](&v.[[ .fieldName ]]), [[ if .hasLocation ]] [[ .location ]], [[ end ]] [[ .displayName ]])
`, gengo.Args{
								"fieldName":      sf.Name,
								"displayName":    strconv.Quote(sf.DisplayName),
								"validateFnName": validateFnName(fnName, NameField, sf.Name),
								"hasLocation":    sf.Location != "",
								"location":       fmt.Sprintf("%s(%s)", d.ReflectTypeLit(reflect.TypeOf(sf.Location)), d.ValueLit(sf.Location)),
							})
						}, []string{})
					},
				})
			}
		},
	})
}

type StructField struct {
	Name        string
	DisplayName string
	Location    validator.Location
	Type        typesutil.Type
	Validator   *validator.ValidatorLoader
}

func walkStruct(s typesutil.Type, each func(sf *StructField), parents []string) {
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fieldName := field.Name()
		fieldTag := field.Tag()

		displayNameTag, hasDisplayNameTag := fieldTag.Lookup("json")
		if !hasDisplayNameTag {
			displayNameTag, hasDisplayNameTag = fieldTag.Lookup("name")
		}

		locationTag, _ := fieldTag.Lookup("in")

		validateTag, _ := fieldTag.Lookup("validate")

		defaultValue, hasDefaultValue := fieldTag.Lookup("default")

		displayName := strings.Split(displayNameTag, ",")[0]
		omitempty := strings.Contains(displayNameTag, "omitempty")

		if !ast.IsExported(fieldName) || displayName == "-" {
			continue
		}

		fieldType := field.Type()

		for {
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			} else {
				break
			}
		}

		if field.Anonymous() {
			switch fieldType.Kind() {
			case reflect.Struct:
				if !hasDisplayNameTag {
					walkStruct(fieldType, each, parents)
					continue
				}
			case reflect.Interface:
				continue
			}
		}

		sf := &StructField{}

		sf.Name = fieldName
		sf.DisplayName = displayName
		sf.Type = field.Type()
		sf.Location = validator.Location(locationTag)

		if validateTag != "-" {
			v := validator.ValidatorMgrDefault.MustCompile(context.Background(), []byte(validateTag), fieldType, func(rule validator.RuleModifier) {
				rule.SetOptional(omitempty)

				if hasDefaultValue {
					rule.SetDefaultValue([]byte(defaultValue))
				}
			})

			vl := v.(*validator.ValidatorLoader)
			if vl != nil && vl.Validator != nil {
				sf.Validator = vl
			}
		}

		each(sf)
	}
}
