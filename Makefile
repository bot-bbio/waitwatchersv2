# Makefile for WaitWatchersV2

.PHONY: all build-wasm build-cli test clean

all: build-wasm build-cli

build-wasm:
	GOOS=js GOARCH=wasm go build -o main.wasm cmd/wasm/main.go

build-cli:
	go build -o mach.exe main.go

test:
	go test ./internal/...

clean:
	rm -f main.wasm mach.exe
