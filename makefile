run:
	go run main.go apiserver

lint:
	golangci-lint run --timeout 10m ./... --fix

tidy:
	go mod tidy && go mod vendor

test:
	go clean -testcache && go test ./...

cover:
	go clean -testcache && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

swagger:
	swag init

docker:
	docker-compose -f ./build/docker-compose.yml up -d

gql:
	go run github.com/99designs/gqlgen generate --config=./config/gqlgen.yml

grpc_client:
	protoc -I app/api/grpc app/api/grpc/*/*.proto \
    --go_out=app/api/grpc --go_opt=paths=source_relative \
    --go-grpc_out=app/api/grpc --go-grpc_opt=paths=source_relative

grpc_server:
	protoc -I internal/grpc internal/grpc/*/*.proto \
    --go_out=internal/grpc --go_opt=paths=source_relative \
    --go-grpc_out=internal/grpc --go-grpc_opt=paths=source_relative
