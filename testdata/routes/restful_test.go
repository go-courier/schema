package routes

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/go-courier/httptransport"
	"github.com/go-courier/httptransport/httpx"
	"github.com/go-courier/httptransport/transformers"
	"github.com/go-courier/schema/testdata/a"
	contextx "github.com/go-courier/x/context"
	"github.com/julienschmidt/httprouter"
)

type GetByIDCopy GetByID

var (
	req = &GetByID{}
)

func init() {
	req.ID = "1z"
	req.Authorization = "Bearer XXX"
	req.Protocol = a.PROTOCOL__HTTP
	req.Label = []string{"label-1", "label-2"}
	req.Name = "name"
}

func toPath() string {
	return transformers.NewPathnamePattern(req.Path()).Stringify(transformers.ParamsFromMap(map[string]string{
		"id": req.ID,
	}))
}

func newRequest() (*http.Request, error) {
	r, err := http.NewRequest(req.Method(), toPath(), nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}

	p, err := req.Protocol.MarshalText()
	if err != nil {
		return nil, err
	}

	query["protocol"] = []string{string(p)}
	query["label"] = req.Label
	query["name"] = []string{req.Name}
	r.Header["Authorization"] = []string{req.Authorization}

	r.URL.RawQuery = query.Encode()

	return r, nil
}

func newIncomingRequest(path string) *http.Request {
	req, _ := newRequest()
	params, _ := transformers.NewPathnamePattern(path).Parse(req.URL.Path)
	return req.WithContext(contextx.WithValue(req.Context(), httprouter.ParamsKey, params))
}

func BenchmarkDecodeFrom(b *testing.B) {
	rtm := httptransport.NewRequestTransformerMgr(nil, nil)

	b.Run("by reflect", func(b *testing.B) {
		r := newIncomingRequest(req.Path())

		req := GetByIDCopy{}

		rt, _ := rtm.NewRequestTransformer(context.Background(), reflect.TypeOf(req))

		for i := 0; i < b.N; i++ {
			_ = rt.DecodeFromRequestInfo(context.Background(), httpx.NewRequestInfo(r), &req)
		}

		b.Log(req)
	})

	b.Run("by generated", func(b *testing.B) {
		r := newIncomingRequest(req.Path())
		req := GetByID{}

		rt, _ := rtm.NewRequestTransformer(context.Background(), reflect.TypeOf(req))

		for i := 0; i < b.N; i++ {
			_ = rt.DecodeFromRequestInfo(context.Background(), httpx.NewRequestInfo(r), &req)
		}

		b.Log(req)
	})
}
