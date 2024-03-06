# Go Backend
- [Go Backend](#go-backend)
- [Project structure](#project-structure)
  - [app folder](#app-folder)
  - [build folder](#build-folder)
  - [cmd folder](#cmd-folder)
  - [config folder](#config-folder)
  - [database folder](#database-folder)
  - [docs folder](#docs-folder)
  - [global folder](#global-folder)
  - [log folder](#log-folder)
  - [model folder](#model-folder)
  - [util folder](#util-folder)
- [Add Restful api](#add-restful-api)
- [Add GraphQL api](#add-graphql-api)
- [Add gRPC api](#add-grpc-api)
- [Start the server](#start-the-server)


# Project structure
```
go_backend/
├─ app/
│  ├─ api/
│  │  ├─ graphql/
│  │  │  ├─ gql/
│  │  │  │  ├─ schema/
│  │  │  │  │  ├─ user.graphqls
│  │  │  │  ├─ model/
│  │  │  │  ├─ generated/
│  │  │  ├─ handler.go
│  │  ├─ grpc/
│  │  │  ├─ user/
│  │  │  │  ├─ user.proto
│  │  ├─ rest/
│  │  │  ├─ v1/
│  │  │  │  ├─ handler/
│  │  │  │  │  ├─ user/
│  │  │  │  ├─ middleware/
│  │  │  │  ├─ response.go
│  │  ├─ router.go
│  ├─ repo/
│  │  ├─ mysql/
│  │  │  ├─ user/
│  │  ├─ redis/
│  │  │  ├─ user/
│  ├─ service/
│  │  ├─ user/
├─ build/
│  ├─ docker-compose.yml
├─ cmd/
│  ├─ apiserver.go
│  ├─ root.go
├─ config/
│  ├─ config.example.yaml
│  ├─ gqlgen.yml
├─ database/
│  ├─ mysql/
│  │  ├─ mysql.go
│  ├─ redis/
│  │  ├─ redis.go
├─ docs/
│  ├─ docs.go
│  ├─ swagger.json
│  ├─ swagger.yaml
├─ global/
│  ├─ global.go
├─ log/
│  ├─ logger.go
├─ model/
│  ├─ user.go
├─ util/
│  ├─ crypto.go
│  ├─ jwt.go
│  ├─ jwt_test.go
│  ├─ util.go
│  ├─ util_test.go
├─ .gitignore
├─ go.mod
├─ go.sum
├─ LICENSE
├─ main.go
├─ makefile
├─ README.md
```
## app folder
* api: contains all the api related code, such as rest, graphql, grpc
* repo: contains all the repository related code
* service: contains all the service related code
* router.go: contains the router setup

## build folder
* docker-compose.yml: contains the docker-compose file

## cmd folder
* apiserver.go: contains the main function
* root.go: contains the root command

## config folder
* config.example.yaml: contains the example config file
* gqlgen.yml: contains the gqlgen config file

## database folder
* mysql: contains the mysql connection setup
* redis: contains the redis connection setup

## docs folder
* docs.go: contains the swagger setup
* swagger.json: contains the swagger json file
* swagger.yaml: contains the swagger yaml file

## global folder
* global.go: contains the global variable

## log folder
* logger.go: contains the logger setup

## model folder
* user.go: contains the user model

## util folder
* contains the util code


# Add Restful api
1. Add route file in app/api/http/routes, for example user.go

# Add GraphQL api
1. Add schema file in app/api/graphql/gql/schema, for example user.graphql
2. Set up the yml file for gqlgen, for example gqlgen.yml in config folder
3. Generate graphql code using gqlgen command as below
```bash
go run github.com/99designs/gqlgen generate --config=./config/gqlgen.yml
```
1. Easy way to generate graphql code using makefile
```bash
make gql
```


# Add gRPC api
1. Add proto file in app/api/grpc, for example user.proto
2. Generate gRPC code from proto file using protoc command as below
```bash
protoc -I app/api/grpc app/api/grpc/*/*.proto \
    --go_out=app/api/grpc --go_opt=paths=source_relative \
    --go-grpc_out=app/api/grpc --go-grpc_opt=paths=source_relative
```
3. Easy way to generate gRPC code using makefile
```bash
make grpc
```

# Start the server
```bash
make docker && make run
```
