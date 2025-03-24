# Use the official Golang image as a build stage
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -o /helpbot

# Use a minimal base image for the final stage
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /helpbot /helpbot

# Command to run
CMD ["/helpbot"]