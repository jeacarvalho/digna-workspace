# Multi-stage Dockerfile for Digna application
# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy all files
COPY . .

# Download dependencies
RUN go mod download

# Build the application with CGO enabled for SQLite
WORKDIR /app/modules/ui_web
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o digna .

# Final stage
FROM alpine:latest

# Install runtime dependencies for SQLite
RUN apk add --no-cache ca-certificates libc6-compat

# Create non-root user
RUN addgroup -S digna && adduser -S digna -G digna

# Create data directory
RUN mkdir -p /var/lib/digna/data && chown -R digna:digna /var/lib/digna

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/modules/ui_web/digna .

# Copy templates and static files
COPY --from=builder /app/modules/ui_web/templates ./templates
COPY --from=builder /app/modules/ui_web/static ./static

# Set ownership
RUN chown -R digna:digna /app

# Switch to non-root user
USER digna

# Expose port
EXPOSE 8090

# Set environment variables with defaults
ENV DIGNA_PORT=8090
ENV DIGNA_DATA_DIR=/var/lib/digna/data
ENV DIGNA_LOG_LEVEL=info

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${DIGNA_PORT}/health || exit 1

# Run the application
CMD ["./digna"]