FROM golang:1.21.3-alpine AS builder
LABEL stage=gobuilder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssl

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

RUN mkdir -p builderCert
RUN openssl genrsa -out builderCert/access 4096
RUN openssl rsa -in builderCert/access -pubout -out builderCert/access.pub

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/builderCert ./cert
COPY --from=builder /app/db/migrations ./db/migrations
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

EXPOSE 8000

CMD ["./main"]