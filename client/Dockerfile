# Use a Go base image
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# # Copy the go mod and sum files
# COPY client/go.mod client/go.sum ./

# # Download the dependencies
# RUN go mod download

# # Copy the source code into the container
# COPY client /app

# # Build the Go app
# RUN go build -o client main.go

# # Start a new stage from scratch
# FROM debian:bullseye-slim

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/client .

# # Command to run the executable
# CMD ["./client"]
