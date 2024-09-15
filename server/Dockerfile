# Use a Go base image
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o server main.go

EXPOSE 50051

CMD ["./server"]

# # Copy the go mod and sum files
# COPY server/go.mod server/go.sum ./

# # Download the dependencies
# RUN go mod download

# # Copy the source code into the container
# COPY server /app

# # Build the Go app
# RUN go build -o server main.go

# # Start a new stage from scratch
# FROM debian:bullseye-slim

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/server .

# # Ensure the binary has execute permissions
# RUN chmod a+x ./server

# # Expose port 50051 to the outside world
# EXPOSE 50051

# # Command to run the executable
# CMD ["./server"]
