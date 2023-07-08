

.PHONY: build

all: build

test:
	@go test -v ./....

build:
	@go build

add:
	@go work use -r modules
	@go work use -r tools
