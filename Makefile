all: build

build:
	go build -v -buildmode=c-shared -o ./bin/libgobtcsuite.so ./src/*.go

.PHONY: all
