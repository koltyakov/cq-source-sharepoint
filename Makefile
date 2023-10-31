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
package: # before using it, craft separate /docs folder and place the command to release pipeline
	go run main.go package --docs-dir docs -m @CHANGELOG.md $(git describe --tags --abbrev=0) .