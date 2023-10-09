# Simple Bank

A simple backend banking app made using Go

# Prerequisites

1. [Go](https://golang.org/doc/install)
2. [PostgreSQL](https://www.postgresql.org/download/)
3. [Docker](https://docs.docker.com/get-docker/)
4. [Golang Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)
5. [Golang air](https://github.com/cosmtrek/air#installation)
6. [Golang Mock](https://github.com/uber-go/mock#installation)

# Setup

1. Clone this repository
2. Run `make postgres` to start the PostgreSQL database
3. Run `make createdb` to create the database
4. Run `make migrateup` to run the database migrations

# Run

1. After [setting up](#setup) the project, run `make devserver` to start the server
    - See [server.go](./src/api/server.go) for the available endpoints

# Tests

1. Run `make test` after [setting up](#setup) the project

# Clean up

1. Run `make migratedown` to clean up the database migrations
2. Run `make dropdb` to drop the database

# Credit

[This](https://www.youtube.com/playlist?list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE)
