postgres:
	 docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	 docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	cd src && migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 

migrateup1:
	cd src && migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	cd src && migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 

migratedown1:
	cd src && migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	 sqlc generate

test:
	cd src && go test -v -cover -test.coverprofile=coverage.cov ./...

mock:
	cd src && mockgen -package mockdb -destination db/mock/store.go github.com/beatzoid/simple-bank/db/sqlc Store

devserver:
	cd src && air

runserver:
	cd src && go run main.go

build:
	cd src && go build -v ./...

dockerbuild:
	 docker build -t simplebank:latest .

dockerrun:
	docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres12:5432/simple_bank?sslmode=disable" simplebank:latest

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test mock devserver runserver build dockerbuild dockerrun