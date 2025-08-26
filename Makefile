lint:
	golangci-lint run

build:
	go build -o bin/gendiff ./cmd/gendiff/main.go	
	
test:
	go test -v ./...
