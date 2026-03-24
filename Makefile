
all: run

build:
	go build -o build/gophermart cmd/gophermart/main.go

run: 
	go run cmd/gophermart/main.go

test:
	go clean -testcache
	go test -count 1 -v -cover ./...

clean:
	rm -f build

coverage:
	go test -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: build run test clean all
