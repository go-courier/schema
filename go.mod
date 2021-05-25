module github.com/go-courier/schema

go 1.16

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/go-courier/courier v1.5.0
	github.com/go-courier/gengo v0.0.0-20210806091446-45a5791f4b63
	github.com/go-courier/httptransport v1.20.5
	github.com/go-courier/reflectx v1.3.4
	github.com/go-courier/statuserror v1.2.1
	github.com/go-courier/validator v1.6.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/onsi/gomega v1.15.0
	github.com/pkg/errors v0.9.1
)

replace github.com/go-courier/httptransport => ../../../github.com/go-courier/httptransport

replace github.com/go-courier/gengo => ../../../github.com/go-courier/gengo
