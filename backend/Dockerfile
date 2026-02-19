# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary (CGO not needed for pgx - pure Go)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /tms ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS and curl for healthcheck
RUN apk add --no-cache ca-certificates curl

# Copy binary from builder
COPY --from=builder /tms /app/tms

# Copy migrations
COPY migrations/ /app/migrations/

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV DATABASE_URL=postgres://postgres:postgres@postgres:5432/tms?sslmode=disable

# Run the binary
CMD ["/app/tms"]
