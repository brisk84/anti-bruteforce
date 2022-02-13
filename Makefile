generate:
	protoc --go_out=internal/server --go_opt=paths=source_relative --go-grpc_out=internal/server --go-grpc_opt=paths=source_relative api/abti-bruteforce.proto

build:
	GOOS=linux GOARCH=amd64  go build -o ./deploy/ab-srv ./cmd/server/
	GOOS=linux GOARCH=amd64  go build -o ./deploy/ab-client ./cmd/client/
	chmod +x ./deploy/ab-srv
	chmod +x ./deploy/ab-client

run:
	./deploy/ab-srv

test:
	go test -race -count 100 ./internal/app
	./deploy/int_test.sh

test_act:
	go test -race -count 100 ./internal/app
	./deploy/int_test_act.sh

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...
