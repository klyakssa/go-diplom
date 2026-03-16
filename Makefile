
all: run

build:
	go build -o build/gophermart cmd/gophermart/main.go

run: 
	go run cmd/gophermart/main.go

test:
	go test ./...

clean:
	rm -f build

.PHONY: build run test clean all
