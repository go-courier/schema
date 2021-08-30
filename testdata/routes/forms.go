package routes

import (
	"context"
	"mime/multipart"

	"github.com/go-courier/courier"
	"github.com/go-courier/httptransport"
	"github.com/go-courier/httptransport/httpx"
)

var FormsRouter = courier.NewRouter(httptransport.Group("/forms"))

func init() {
	RootRouter.Register(FormsRouter)

	FormsRouter.Register(courier.NewRouter(FormURLEncoded{}))
	FormsRouter.Register(courier.NewRouter(FormMultipartWithFile{}))
	FormsRouter.Register(courier.NewRouter(FormMultipartWithFiles{}))
}

// Form URL Encoded
type FormURLEncoded struct {
	httpx.MethodPost `path:"/urlencoded"`
	FormData         struct {
		String string   `name:"string"`
		Slice  []string `name:"slice"`
		Data   Data     `name:"data"`
	} `in:"body" mime:"urlencoded"`
}

func (req FormURLEncoded) Output(ctx context.Context) (resp interface{}, err error) {
	return
}

// Form Multipart
type FormMultipartWithFile struct {
	httpx.MethodPost `path:"/multipart"`

	FormData struct {
		// @deprecated
		String string                `name:"string,omitempty"`
		Slice  []string              `name:"slice,omitempty"`
		Data   Data                  `name:"data,omitempty"`
		File   *multipart.FileHeader `name:"file" validate:"-"`
	} `in:"body" mime:"multipart"`
}

func (req FormMultipartWithFile) Output(ctx context.Context) (resp interface{}, err error) {
	return
}

// Form Multipart With Files
type FormMultipartWithFiles struct {
	httpx.MethodPost `path:"/multipart-with-files"`

	FormData struct {
		Files []*multipart.FileHeader `name:"files" validate:"-"`
	} `in:"body" mime:"multipart"`
}

func (FormMultipartWithFiles) Output(ctx context.Context) (resp interface{}, err error) {
	return
}
