all: build

build:
	go build -v -buildmode=c-shared -o ./bin/linux/libgobtcsuite.so ./src/*.go

.PHONY: all
