# Start from golang base image to build the server
FROM golang:1.19-alpine as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Tools need to compile drivers
RUN apk update && apk add --no-cache make gcc g++ musl-dev binutils autoconf automake libtool pkgconfig check-dev file patch

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
#RUN go build -a -o nw .
RUN go build -a --ldflags '-linkmode external -extldflags "-static"' -o nw .



# Start a new stage from scratch
FROM alpine:3.16

# Add Maintainer info
LABEL maintainer="Jocelyn GENNESSEAUX"

# Certificates
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

# Copy the Pre-built binary file from the previous stage.
# Don't forget to copy the .env file
COPY --from=builder /build/nw /nw/nw

WORKDIR /nw

# Command to run the executable
CMD ["/nw/nw"]
