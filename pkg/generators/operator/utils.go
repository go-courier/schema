package openapi

import (
	"github.com/go-courier/courier"
	"github.com/go-courier/httptransport"
	"github.com/go-courier/httptransport/httpx"
	"github.com/go-courier/statuserror"
	"go/ast"
	"go/types"
	"reflect"
	"strings"
)

const (
	XStatusErrs = `x-status-errors`
)

var (
	pkgImportPathHttpTransport = importGoPath(reflect.TypeOf(httptransport.HttpRouteMeta{}).PkgPath())
	pkgImportPathStatusError   = importGoPath(reflect.TypeOf(statuserror.StatusErr{}).PkgPath())
	pkgImportPathHttpx         = importGoPath(reflect.TypeOf(httpx.Response{}).PkgPath())
	pkgImportPathCourier       = importGoPath(reflect.TypeOf(courier.Router{}).PkgPath())
)

func importGoPath(importPath string) string {
	parts := strings.Split(importPath, "/vendor/")
	return parts[len(parts)-1]
}

func isRouterType(typ types.Type) bool {
	return strings.HasSuffix(typ.String(), pkgImportPathCourier+".Router")
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

func isFromHttpTransport(typ types.Type) bool {
	return strings.Contains(typ.String(), pkgImportPathHttpTransport+".")
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
