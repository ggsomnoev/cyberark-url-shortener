FROM golang:1.24 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /build

COPY go.mod go.sum ./
COPY vendor ./vendor
ENV GOFLAGS=-mod=vendor

COPY internal ./internal
COPY cmd/api ./cmd/api

COPY migrations ./migrations

RUN go build -o api ./cmd/api

FROM alpine:latest

WORKDIR /app
COPY --from=builder /build/api .
COPY --from=builder /build/migrations ./migrations

CMD ["./api"]