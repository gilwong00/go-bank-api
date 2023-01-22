postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=postgres --owner=postgres bank_api

dropdb:
	docker exec -it postgres14 dropdb bank_api

migrateup:
	migrate -path "pkg/db/migrations" -database "postgres://postgres:postgres@localhost:5432/bank_api?sslmode=disable" up

migratedown:
	migrate -path "pkg/db/migrations" -database "postgres://postgres:postgres@localhost:5432/bank_api?sslmode=disable" down

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover ./...

mockDB:
	mockgen -package mockDB -destination pkg/db/mock/store.go go-bank-api/sqlc Store

.PHONY: postgres createdb dropdb sqlc server test
