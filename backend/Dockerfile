# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /tms ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS, sqlite, and wget for healthcheck
RUN apk add --no-cache ca-certificates sqlite-libs wget

# Copy binary from builder
COPY --from=builder /tms /app/tms

# Copy migrations
COPY migrations/ /app/migrations/

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV DB_PATH=/app/data/tms.db

# Run the binary
CMD ["/app/tms"]
