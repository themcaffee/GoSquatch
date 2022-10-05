# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as build

WORKDIR /

# Copy the Go Modules manifests and get dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Build application
COPY *.go ./
RUN go build -o Squatch

## Build a small image

FROM alpine:latest

WORKDIR /

COPY --from=build /Squatch /Squatch
COPY $INPUT_SRC_DIR /$INPUT_SRC_DIR

ENTRYPOINT ["./Squatch"]