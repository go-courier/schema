package routes

import (
	"context"
	"net/http"

	"github.com/go-courier/courier"
	"github.com/go-courier/httptransport"
	"github.com/go-courier/httptransport/httpx"
	"github.com/go-courier/schema/testdata/a"
	"github.com/go-courier/schema/testdata/b"
	"github.com/go-courier/statuserror"
	perrors "github.com/pkg/errors"
)

var RestfulRouter = courier.NewRouter(httptransport.Group("/restful"), &MustValidAccount{})

func init() {
	RootRouter.Register(RestfulRouter)

	RestfulRouter.Register(courier.NewRouter(HealthCheck{}))
	RestfulRouter.Register(courier.NewRouter(Create{}))
	RestfulRouter.Register(courier.NewRouter(DataProvider{}, UpdateByID{}))
	RestfulRouter.Register(courier.NewRouter(DataProvider{}, GetByID{}))
	RestfulRouter.Register(courier.NewRouter(DataProvider{}, RemoveByID{}))
}

type HealthCheck struct {
	httpx.MethodHead

	PullPolicy b.PullPolicy `name:"pullPolicy,omitempty" in:"query"`
}

func (HealthCheck) Output(ctx context.Context) (interface{}, error) {
	return nil, nil
}

// Create
type Create struct {
	httpx.MethodPost
	Data Data `in:"body"`
}

func (req Create) Output(ctx context.Context) (interface{}, error) {
	return &req.Data, nil
}

type Data struct {
	ID        string     `json:"id"`
	Label     string     `json:"label"`
	PtrString *string    `json:"ptrString,omitempty"`
	Data      *Data      `json:"data,omitempty"`
	SubData   *SubData   `json:"subData,omitempty"`
	Protocol  a.Protocol `json:"protocol,omitempty"`
}

type SubData struct {
	Name string `json:"name"`
}

// get by id
type GetByID struct {
	httpx.MethodGet
	Protocol a.Protocol `name:"protocol,omitempty" in:"query"`
	Name     string     `name:"name,omitempty" in:"query"`
	Label    []string   `name:"label,omitempty" in:"query"`
}

func (req GetByID) Output(ctx context.Context) (interface{}, error) {
	data := DataFromContext(ctx)
	if len(req.Label) > 0 {
		data.Label = req.Label[0]
	}
	return data, nil
}

// remove by id
type RemoveByID struct {
	httpx.MethodDelete
}

func callWithErr() error {
	return statuserror.Wrap(perrors.New("test"), http.StatusInternalServerError, "InternalServerX")
}

func (RemoveByID) Output(ctx context.Context) (interface{}, error) {
	if false {
		return nil, callWithErr()
	}
	return nil, statuserror.Wrap(perrors.New("test"), http.StatusInternalServerError, "InternalServer")
}

// update by id
type UpdateByID struct {
	httpx.MethodPut
	Data Data `in:"body"`
}

func (req UpdateByID) Output(ctx context.Context) (interface{}, error) {
	return nil, perrors.Errorf("something wrong")
}

type DataProvider struct {
	httpx.Method `path:"/:id"`

	ID string `name:"id" in:"path" validate:"@string[6,]"`
}

func (DataProvider) ContextKey() string {
	return "DataProvider"
}

func DataFromContext(ctx context.Context) *Data {
	return ctx.Value("DataProvider").(*Data)
}

func (req DataProvider) Output(ctx context.Context) (interface{}, error) {
	return &Data{
		ID: req.ID,
	}, nil
}
