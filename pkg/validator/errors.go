package validator

import (
	"bytes"
	"container/list"
	"fmt"
	"reflect"

	"github.com/go-courier/reflectx"
)

func NewErrorSet(paths ...string) *ErrorSet {
	return &ErrorSet{
		root:   paths,
		errors: list.New(),
	}
}

type ErrorSet struct {
	root   []string
	errors *list.List
}

func (errorSet *ErrorSet) AddErr(err error, keyPathNodes ...interface{}) {
	if err == nil {
		return
	}

	errorSet.errors.PushBack(&FieldError{
		Field: KeyPath(keyPathNodes),
		Error: err,
	})
}

func (errorSet *ErrorSet) Each(cb func(fieldErr *FieldError)) {
	l := errorSet.errors

	for e := l.Front(); e != nil; e = e.Next() {
		if fieldErr, ok := e.Value.(*FieldError); ok {
			cb(fieldErr)
		}
	}
}

func (errorSet *ErrorSet) Flatten() *ErrorSet {
	set := NewErrorSet(errorSet.root...)

	errorSet.Each(func(fieldErr *FieldError) {
		if subSet, ok := fieldErr.Error.(*ErrorSet); ok {
			subSet.Flatten().Each(func(subSetFieldErr *FieldError) {
				set.AddErr(subSetFieldErr.Error, append(fieldErr.Field, subSetFieldErr.Field...)...)
			})
		} else {
			set.AddErr(fieldErr.Error, fieldErr.Field...)
		}
	})

	return set
}

func (errorSet *ErrorSet) Len() int {
	return errorSet.Flatten().errors.Len()
}

func (errorSet *ErrorSet) Err() error {
	if errorSet.errors.Len() == 0 {
		return nil
	}
	return errorSet
}

func (errorSet *ErrorSet) Error() string {
	set := errorSet.Flatten()

	buf := bytes.Buffer{}

	set.Each(func(fieldErr *FieldError) {
		buf.WriteString(fmt.Sprintf("%s %s", fieldErr.Field, fieldErr.Error))
		buf.WriteRune('\n')
	})

	return buf.String()
}

type FieldError struct {
	Field KeyPath
	Error error `json:"msg"`
}

type KeyPath []interface{}

func (keyPath KeyPath) String() string {
	buf := &bytes.Buffer{}
	for i := 0; i < len(keyPath); i++ {
		switch keyOrIndex := keyPath[i].(type) {
		case string:
			if buf.Len() > 0 {
				buf.WriteRune('.')
			}
			buf.WriteString(keyOrIndex)
		case int:
			buf.WriteString(fmt.Sprintf("[%d]", keyOrIndex))
		}
	}
	return buf.String()
}

type MissingRequired struct{}

func (MissingRequired) Error() string {
	return "missing required field"
}

type NotMatchError struct {
	Target  string
	Current interface{}
	Pattern string
}

func (err *NotMatchError) Error() string {
	return fmt.Sprintf("%s %s not match %v", err.Target, err.Pattern, err.Current)
}

type MultipleOfError struct {
	Target     string
	Current    interface{}
	MultipleOf interface{}
}

func (e *MultipleOfError) Error() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(e.Target)
	buf.WriteString(fmt.Sprintf(" should be multiple of %v", e.MultipleOf))
	buf.WriteString(fmt.Sprintf(", but got invalid value %v", e.Current))
	return buf.String()
}

type NotInEnumError struct {
	Target  string
	Current interface{}
	Enums   []interface{}
}

func (e *NotInEnumError) Error() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(e.Target)
	buf.WriteString(" should be one of ")

	for i, v := range e.Enums {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%v", v))
	}

	buf.WriteString(fmt.Sprintf(", but got invalid value %v", e.Current))

	return buf.String()
}

type OutOfRangeError struct {
	Target           string
	Current          interface{}
	Minimum          interface{}
	Maximum          interface{}
	ExclusiveMaximum bool
	ExclusiveMinimum bool
}

func (e *OutOfRangeError) Error() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(e.Target)
	buf.WriteString(" should be")

	if e.Minimum != nil {
		buf.WriteString(" larger")
		if e.ExclusiveMinimum {
			buf.WriteString(" or equal")
		}

		buf.WriteString(fmt.Sprintf(" than %v", reflectx.Indirect(reflect.ValueOf(e.Minimum)).Interface()))
	}

	if e.Maximum != nil {
		if e.Minimum != nil {
			buf.WriteString(" and")
		}

		buf.WriteString(" less")
		if e.ExclusiveMaximum {
			buf.WriteString(" or equal")
		}

		buf.WriteString(fmt.Sprintf(" than %v", reflectx.Indirect(reflect.ValueOf(e.Maximum)).Interface()))
	}

	buf.WriteString(fmt.Sprintf(", but got invalid value %v", e.Current))

	return buf.String()
}

func NewUnsupportedTypeError(typ string, rule string, msgs ...string) *UnsupportedTypeError {
	return &UnsupportedTypeError{
		rule: rule,
		typ:  typ,
		msgs: msgs,
	}
}

type UnsupportedTypeError struct {
	msgs []string
	rule string
	typ  string
}

func (e *UnsupportedTypeError) Error() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(e.rule)
	buf.WriteString(" could not validate type ")
	buf.WriteString(e.typ)

	for i, msg := range e.msgs {
		if i == 0 {
			buf.WriteString(": ")
		} else {
			buf.WriteString("; ")
		}
		buf.WriteString(msg)
	}

	return buf.String()
}
