# 🏗️ Stage 1: Build the Golang Application
FROM golang:1.23 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4 && swag init -g cmd/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/user-api main.go

# 🏁 Stage 2: Create a minimal production image
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set working directory inside the final container
WORKDIR /root

# Copy the built binary from the builder stage
COPY --from=builder /app/user-api /root/user-api

# Copy Swagger documentation
COPY --from=builder /app/docs /root/docs

# Ensure the binary has execute permissions
RUN chmod +x /root/user-api

# Expose API port
EXPOSE 8080

# Run the API
CMD ["/root/user-api"]