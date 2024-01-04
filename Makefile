.PHONY: check
check: fmt vet lint test

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: vet
vet:
	@go vet ./...

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@go test -v -coverprofile=coverage.txt -covermode=atomic ./...
