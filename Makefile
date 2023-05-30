DB_URL=postgres://postgres:postgres@localhost:5432/bank_api?sslmode=disable

postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=postgres --owner=postgres bank_api

dropdb:
	docker exec -it postgres14 dropdb bank_api

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

.PHONY: postgres createdb dropdb sqlc server test migrateuplatest migratedownlast mock
