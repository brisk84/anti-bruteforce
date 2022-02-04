generate:
	protoc --go_out=internal/server --go_opt=paths=source_relative --go-grpc_out=internal/server --go-grpc_opt=paths=source_relative api/abti-bruteforce.proto

build:
	go build -o deploy/ab-srv cmd/main.go
	go build -o deploy/ab-client cmd/client/client.go
	chmod +x deploy/ab-srv
	chmod +x deploy/ab-client

run:
	./deploy/ab-srv

test:
	go test ./internal/app
	cd deploy; ./int_test.sh