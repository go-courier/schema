package openapi

import (
	"go/ast"
	"go/constant"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-courier/httptransport/httpx"
	"github.com/go-courier/statuserror"
)

const (
	XStatusErrs = `x-status-errors`
)

var (
	pkgImportPathStatusError = importGoPath(reflect.TypeOf(statuserror.StatusErr{}).PkgPath())
	pkgImportPathHttpx       = importGoPath(reflect.TypeOf(httpx.Response{}).PkgPath())
)

func importGoPath(importPath string) string {
	parts := strings.Split(importPath, "/vendor/")
	return parts[len(parts)-1]
}

func isNil(typ types.Type) bool {
	return typ == nil || typ.String() == types.Typ[types.UntypedNil].String()
}

func isHttpxResponse(typ types.Type) bool {
	return strings.HasSuffix(typ.String(), pkgImportPathHttpx+".Response")
}

func isStatusErr(typ types.Type) bool {
	return strings.HasSuffix(typ.String(), pkgImportPathStatusError+".StatusErr")
}

func getIdentChainOfCallFunc(expr ast.Expr) (list []*ast.Ident) {
	switch e := expr.(type) {
	case *ast.CallExpr:
		list = append(list, getIdentChainOfCallFunc(e.Fun)...)
	case *ast.SelectorExpr:
		list = append(list, getIdentChainOfCallFunc(e.X)...)
		list = append(list, e.Sel)
	case *ast.Ident:
		list = append(list, expr.(*ast.Ident))
	}
	return
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
