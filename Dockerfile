# Build stage
FROM golang:1.21.3-alpine3.18 AS builder
WORKDIR /app
COPY . .

WORKDIR /app/src
RUN go build -o ../main ./main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY ./src/app.env .

EXPOSE 8080
CMD [ "/app/main" ]
