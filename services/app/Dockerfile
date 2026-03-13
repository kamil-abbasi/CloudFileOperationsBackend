FROM golang:1.25-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/server/main.go

RUN go build -o healthcheck ./healthcheck.go

FROM golang:1.25-alpine AS dev

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY --from=builder /usr/src/app/healthcheck .

CMD ["air"]

## -----------------
## Production Stage
## -----------------
FROM alpine:latest AS prod

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/healthcheck .

RUN adduser -D appuser
USER appuser

CMD ["./main"]