generate:
	protoc --go_out=internal/server --go_opt=paths=source_relative --go-grpc_out=internal/server --go-grpc_opt=paths=source_relative api/abti-bruteforce.proto

build:
	go version
	go build -o ./deploy/ab-client ./cmd/client/
	go build -o ./deploy/ab-srv ./cmd/server/
	chmod +x ./deploy/ab-client
	chmod +x ./deploy/ab-srv

run:
	./deploy/ab-srv

test:
	go test ./internal/app
	cd deploy; ./int_test.sh

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...
