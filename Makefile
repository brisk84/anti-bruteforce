generate:
	protoc --go_out=internal/server --go_opt=paths=source_relative --go-grpc_out=internal/server --go-grpc_opt=paths=source_relative api/abti-bruteforce.proto

build:
	go build -o ab-srv ./cmd/server/
	go build -o ab-client ./cmd/client/
	chmod +x ./ab-srv
	chmod +x ./ab-client

run:
	go run ./cmd/server

test:
	go test -race -count 100 ./internal/app

test_integration: build
	./deploy/int_test.sh

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...
