# Simple Bank

A simple bank backend application written in Go

# Prerequisites

1. [Go](https://golang.org/doc/install)
2. [PostgreSQL](https://www.postgresql.org/download/)
3. [Docker](https://docs.docker.com/get-docker/)
4. [Golang Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)

# How to run

1. Clone this repository
2. Run `make postgres` to start the PostgreSQL database
3. Run `make createdb` to create the database
4. Run `make migrateup` to run the database migrations

# Clean up

1. Run `make migratedown` to clean up the database migrations
2. Run `make dropdb` to drop the database

# Credit

[This](https://www.youtube.com/playlist?list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE)
