all: run

run: build
	@./bin/author

.PHONY: build
build:
	@go build -o bin/author ./main.go
