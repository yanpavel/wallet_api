FROM golang:1.24.2-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api  && chmod +x api

EXPOSE 8080

CMD ["./api"]