all: build run clean

.PHONY: test
test: lint
	go test $(shell go list ./... | grep -v /vendor)

.PHONY: lint
lint:
	golangci-lint run

build:
	go build -o out

run:
	./out $(port)

clean:
	rm out
