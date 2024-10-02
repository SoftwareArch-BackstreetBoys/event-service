# Use the Golang image as the builder
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install gRPC and Protobuf tools
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy the rest of the application source code
COPY . .

# Run go mod tidy to clean up dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o server main.go

# Expose the gRPC port
EXPOSE 50051

# Command to run the application
CMD ["./server"]
