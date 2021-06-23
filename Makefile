all: build test check
modules:
	go mod tidy
build:
	go build -o ./bin/balanceservice ./cmd/main.go
test:
	go test ./...
check:
	golangci-lint run
generate-grpc:
		protoc -I/usr/local/include -I. \
			-I${GOPATH}/src \
			-I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
            -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 \
			--grpc-gateway_out=logtostderr=true:./api \
			--swagger_out=allow_merge=true,merge_file_name=api:./api \
			--go_out=plugins=grpc:./api \
			./api/api.proto