
FROM golang:1.22.4 AS builder


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o todo_server cmd/main.go


CMD ["./todo_server"]