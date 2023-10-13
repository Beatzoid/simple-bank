# Simple Bank

A simple backend banking app made using Go

# Prerequisites

1. [Go](https://golang.org/doc/install)
2. [PostgreSQL](https://www.postgresql.org/download/)
3. [Docker](https://docs.docker.com/get-docker/)
4. [Golang Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)
5. [Golang air](https://github.com/cosmtrek/air#installation)
6. [Golang Mock](https://github.com/uber-go/mock#installation)

# Setup and Run

1. Clone this repository
2. Run `docker compose up` to initialize and start everything
    - See [server.go](./src/api/server.go) for the available endpoints

# Tests

1. Run `make test`

# Clean up

1. Run `docker compose down` to stop and remove everything
2. Run `docker rmi simplebank-api` to remove the Docker image

# Credit

[This](https://www.youtube.com/playlist?list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE)
