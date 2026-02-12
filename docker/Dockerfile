FROM golang:1.25.1 AS builder

WORKDIR /plugin

COPY go.mod ./

RUN go mod download

COPY . ./

RUN go build -o /sentry

EXPOSE 8080
CMD ["/sentry"]
