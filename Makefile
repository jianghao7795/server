.PHONY: build dev start clean tool help

all: build

build:
	 GOEXPERIMENT=jsonv2 CGO_ENABLED=0 GOOS=linux go build -tags=jsoniter -trimpath -o bin/fiber -ldflags="-s -w" cmd/main.go 

start:
	air

dev:
	go run cmd/main.go -c ./conf/

local:
	go run . -c ./conf/

tool:
	go tool vet . |& grep -v vendor; true
	gofmt -w .

clean:
	go clean -i .

air:
	air -c .air.toml

test:
	go test -v ./...

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make dev: go run ."
	@echo "make clean: remove object files and cached files"
