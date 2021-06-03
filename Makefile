all: build test check
modules:
	go mod tidy
build:
	go build -o ./bin/balanceservice ./cmd/main.go
test:
	go test ./...
check:
	golangci-lint run