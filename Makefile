.PHONY: build test release

build:
	go build -o relic .

test:
	go test ./...

release:
	goreleaser release --clean
