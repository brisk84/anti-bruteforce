generate:
	protoc --go_out=internal/server --go_opt=paths=source_relative --go-grpc_out=internal/server --go-grpc_opt=paths=source_relative api/abti-bruteforce.proto

build:
	go version
	go mod tidy
	go build -o ./deploy/ab-srv ./cmd/server/
	go build -o ./deploy/ab-client ./cmd/client/
	chmod +x ./deploy/ab-srv
	chmod +x ./deploy/ab-client

run:
	./deploy/ab-srv

test:
	go test ./internal/app
	cd deploy; ./int_test.sh

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...
