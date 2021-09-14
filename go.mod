module github.com/go-courier/schema

go 1.17

replace github.com/go-courier/sqlx/v2 => ../sqlx

replace github.com/go-courier/gengo => ../gengo

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/go-courier/courier v1.5.0
	github.com/go-courier/gengo v0.0.0-20210830081703-2ff0a49b8aa4
	github.com/go-courier/httptransport v1.21.6
	github.com/go-courier/sqlx/v2 v2.0.0-00010101000000-000000000000
	github.com/go-courier/statuserror v1.2.1
	github.com/go-courier/x v0.0.11
	github.com/gorilla/handlers v1.5.1
	github.com/onsi/gomega v1.16.0
	github.com/pkg/errors v0.9.1
)

require (
	github.com/fatih/color v1.12.0 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/go-courier/logr v0.0.2 // indirect
	github.com/go-courier/metax v1.3.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/mod v0.5.0 // indirect
	golang.org/x/net v0.0.0-20210913180222-943fd674d43e // indirect
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.5 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
