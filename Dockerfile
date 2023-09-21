# Start from golang base image to build the server
FROM golang:1.21-alpine as builder

# Tools needed to compile
RUN apk update && apk add --no-cache git make gcc g++ musl-dev binutils autoconf automake libtool pkgconfig check-dev file patch

# Set the current working directory inside the container
WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
# Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . ./

# Build env
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOAMD64=v3

# Build the Go app
RUN make build


# Start a new stage from scratch
FROM alpine:3.16

# Add Maintainer info
LABEL maintainer="Jocelyn GENNESSEAUX"

# Certificates
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

# Define working dir
WORKDIR /NotionWatcher

# Copy the Pre-built binary file from the previous stage.
COPY --from=builder /build/build/NotionWatcher /NotionWatcher/notionwatcher

# declare the volume to store the watchers
VOLUME [ "/NotionWatcher/watchers" ]

# declare the volume to store the list of users
VOLUME [ "/NotionWatcher/data" ]

# declare the volume to store logs
VOLUME [ "/NotionWatcher/logs" ]

# Command to run the executable
CMD ["/NotionWatcher/notionwatcher", "watch"]
