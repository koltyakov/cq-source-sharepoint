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

.PHONY: package
package:
	go run main.go package --docs-dir docs -m @CHANGELOG.md v2.0.0 .