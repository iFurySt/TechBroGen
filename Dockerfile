# Build stage
FROM golang:1.24.1 AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM ubuntu:22.04

# Install ca-certificates for HTTPS requests and update packages
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /opt/tbg

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy static assets
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/index.html .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
