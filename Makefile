install:
	go get -u ./... && go mod tidy

.PHONY: test
test:
	go test -timeout 3m ./...

.PHONY: lint
lint:
	@golangci-lint run --timeout 10m

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: build
build:
	go build -o bin/cq-source-sharepoint -v