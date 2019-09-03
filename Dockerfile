# Build stage
FROM golang:1.12.9-alpine3.10 as builder

ENV CGO_ENABLED 0

RUN apk update && apk add git bash
SHELL ["/bin/bash", "-c"]

WORKDIR /src/web-crawler

COPY go.mod go.sum ./
RUN go mod download

# Build the application
COPY . .
RUN go build -a --installsuffix cgo

# Run tests
RUN go test -v ./...

# Production image stage
FROM alpine:3.10

RUN apk --no-cache --update upgrade \
    && apk --no-cache add ca-certificates

COPY --from=builder /src/web-crawler/web-crawler /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/web-crawler"]
