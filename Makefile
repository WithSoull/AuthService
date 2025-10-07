include .env

LOCAL_BIN=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/grpc/
	go get -u github.com/joho/godotenv
	go get -u github.com/WithSoull/platform_common
	go get -u github.com/WithSoull/UserServer
	go get -u github.com/golang-jwt/jwt/v5
	go get -u github.com/gomodule/redigo
	go get -u github.com/gomodule/redigo/redis

generate-api:
	make generate-api-auth

generate-api-auth:
	mkdir -p pkg/auth/v1
	protoc --proto_path api/auth/v1 \
	--go_out=pkg/auth/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/auth/v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/auth/v1/auth.proto


rebuild:
	docker compose down
	docker compose build --no-cache
	docker compose up -d
	docker compose logs -f
