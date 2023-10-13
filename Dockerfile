# Build stage
FROM golang:1.21.3-alpine3.18 AS builder
WORKDIR /app

# Copy source files
COPY . .

# CD into src directory where the code is
WORKDIR /app/src

# Cache the downloaded dependency in the layer
RUN go mod download

# Build the project into a single executible
RUN go build -o ../main ./main.go
# Download the migration tool
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.18
WORKDIR /app
# Copy the executible from the builder stage
COPY --from=builder /app/main .
# Copy the migration tool from the builder stage
COPY --from=builder /app/src/migrate ./migrate
# Copy the env file
COPY ./src/app.env .
# Copy the migration files
COPY ./src/db/migration ./migration
# Copy the wait-for script
COPY ./wait-for.sh .
# Copy the start script
COPY ./start.sh .

# Port that the server runs on
EXPOSE 8080
# CMD is the arguments passed to ENTRYPOINT
# In this case, /app/main is passed to /app/start.sh
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
