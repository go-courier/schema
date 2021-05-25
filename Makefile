t:
	go test -v ./generators

test: tidy
	go test -v -race -failfast ./...

cover: tidy
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

debug:
	go run ./testdata

tidy:
	go mod tidy

fmt:
	goimports -l -w .

dep:
	go get -u ./...