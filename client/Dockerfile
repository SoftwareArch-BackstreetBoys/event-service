# Use the Golang image as the builder
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Run go mod tidy to clean up dependencies
RUN go mod tidy

# Build the Go application (client binary)
RUN go build -o client main.go

# Expose the gRPC port
EXPOSE 6001

# Command to run the application
CMD ["./client"]
