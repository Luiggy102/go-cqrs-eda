ARG GO_VERSION=1.22.7

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY events events
COPY repository repository
COPY search search
COPY database database
COPY models models
COPY feed-service feed-service

RUN go install ./...

# ---------------------------------------------
FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=builder /go/bin .