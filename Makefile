postgres:
	 docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	 docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	cd src && migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 

migratedown:
	cd src && migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 

sqlc:
	cd src && sqlc generate

test:
	cd src && go test -v -cover -test.coverprofile=coverage.cov ./...

server:
	cd src && go run main.go

build:
	cd src && go build -v ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server build