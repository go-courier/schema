t:
	go test -v ./generators

test: tidy
	go test -race ./...

cover: tidy
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

debug:
	go run ./testdata

tidy: fmt
	go mod tidy

fmt:
	goimports -l -w .