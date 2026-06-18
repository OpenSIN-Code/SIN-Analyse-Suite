BINARY = sin-analyse
CGO_ENABLED = 0

.PHONY: build test lint clean check

build:
	CGO_ENABLED=$(CGO_ENABLED) go build -o $(BINARY) ./cmd/$(BINARY)

test:
	CGO_ENABLED=$(CGO_ENABLED) go test ./... -race -count=1 -coverprofile=coverage.out

lint:
	golangci-lint run

check: build lint test

clean:
	rm -f $(BINARY) coverage.out
