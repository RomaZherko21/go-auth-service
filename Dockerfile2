FROM golang:alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

# RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine

# RUN apk update --no-cache && apk add --no-cache ca-certificates

# COPY --from=builder /usr/share/zoneinfo/America/New_York /usr/share/zoneinfo/America/New_York

# ENV TZ GMT+0

WORKDIR /build

COPY --from=builder /build/main /build/main

CMD [". /main"]