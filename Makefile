test:
	go test -race ./...

cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

debug:
	go run ./testdata

fmt:
	goimports -l -w .