package main

import (
	"context"
	"github.com/go-courier/courier"
	"github.com/go-courier/httptransport"
	"github.com/go-courier/reflectx/typesutil"
	"github.com/go-courier/schema/pkg/validator"
	"github.com/go-courier/schema/testdata/routes"
	oldvalidator "github.com/go-courier/validator"
	"github.com/gorilla/handlers"
	"net/http"
)

type ValidatorMgr struct {
	validator.ValidatorMgr
}

func (v *ValidatorMgr) Compile(ctx context.Context, bytes []byte, t typesutil.Type, processors ...oldvalidator.RuleProcessor) (oldvalidator.Validator, error) {
	newProcessors := make([]validator.RuleProcessor, len(processors))

	for i := range newProcessors {
		oldProcess := processors[i]

		newProcessors[i] = func(rule validator.RuleModifier) {
			oldProcess(rule)
		}
	}

	return v.ValidatorMgr.Compile(ctx, bytes, t, newProcessors...)
}

func main() {
	ht := &httptransport.HttpTransport{}
	ht.ValidatorMgr = &ValidatorMgr{ValidatorMgr: validator.ValidatorMgrDefault}
	ht.Middlewares = []httptransport.HttpMiddleware{
		DefaultCORS(),
	}
	ht.SetDefaults()
	courier.Run(routes.RootRouter, ht)
}

var (
	defaultCorsMethods = []string{"GET", "HEAD", "POST"}
	defaultCorsHeaders = []string{"Accept", "Accept-Language", "Content-Language", "Origin"}
)

func DefaultCORS() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods(append(defaultCorsMethods, []string{
			http.MethodConnect,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		}...)),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders(append(defaultCorsHeaders, []string{
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
			"Content-Type",
			"Authorization",
			"X-Request-Id",
			"User-Agent",
		}...)),
		handlers.ExposedHeaders([]string{
			"Content-Type",
			"Origin",
			"b3",
			"User-Agent",
			"X-Requested-With",
			"X-Request-Id",
			"X-Meta",
			// follow https://developer.github.com/v3/rate_limit/
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		}),
		handlers.OptionStatusCode(http.StatusNoContent),
	)
}
