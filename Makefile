test:
	go test -v -race ./...

cover:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...


fmt:
	goimports -l -w .
	gofmt -l -w .