# Start from the latest golang base image
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for go mod download if needed
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o livescore ./cmd/server

# Start a new stage from scratch
FROM alpine:latest
WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/livescore .

# Expose port 3000
EXPOSE 3000

# Command to run
CMD ["./livescore"] 
