FROM golang:1.21.3-alpine AS builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN ln -s /go/bin/linux_amd64/migrate /usr/local/bin/migrate

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# need to make cert after COPY . . but it will make a new cert after every changes 
RUN apk add openssl
RUN mkdir -p cert
RUN openssl genrsa -out cert/access 4096
RUN openssl rsa -in cert/access -pubout -out cert/access.pub

RUN go build -o main ./cmd/main.go

EXPOSE 8000

CMD ["./main"]