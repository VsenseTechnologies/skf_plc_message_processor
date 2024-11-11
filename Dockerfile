#build stage
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/main

#run stage
FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app

COPY --from=build /app/bin/main .

CMD ["./main"]
