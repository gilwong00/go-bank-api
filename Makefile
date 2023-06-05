DB_URL=postgres://postgres:postgres@localhost:5432/bank_db?sslmode=disable

postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
	docker exec -it bank_api_pg createdb --username=postgres --owner=postgres bank_db

dropdb:
	docker exec -it bank_api_pg dropdb bank_db

migration:
	migrate create -ext sql -dir pkg/db/migrations

migrateup:
	migrate -path "pkg/db/migrations" -database "$(DB_URL)" up

migrateuplatest:
	migrate -path "pkg/db/migrations" -database "$(DB_URL)" up 1

migratedown:
	migrate -path "pkg/db/migrations" -database "$(DB_URL)" down

migratedownlast:
	migrate -path "pkg/db/migrations" -database "$(DB_URL)" down 1

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination pkg/db/mock/store.go go-bank-api/pkg/db/sqlc Store

docker_build_image:
	docker build -t bankapi:latest .

build_network:
	docker run --name bankapi -p 5000:5000 -e GIN_MODE=release -e DB_SOURCE="postgres://postgres:postgres@bank_api_pg:5432/bank_api?sslmode=disable" bankapi:latest

proto:
	rm -f rpc/*.go
	protoc --proto_path=proto --go_out=rpc --go_opt=paths=source_relative \
	--go-grpc_out=rpc --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=rpc  --grpc-gateway_opt paths=source_relative \
	proto/*.proto

evans:
	evans  --host localhost --port 6000 -r repl

.PHONY: postgres createdb dropdb sqlc server test migrateuplatest migratedownlast mock docker_build_image proto evans
