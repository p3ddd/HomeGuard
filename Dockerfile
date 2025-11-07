# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o homeguard .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/homeguard .

# Copy example config
COPY config.example.yaml .

# Create a non-root user
RUN addgroup -g 1000 homeguard && \
    adduser -D -u 1000 -G homeguard homeguard && \
    chown -R homeguard:homeguard /app

USER homeguard

EXPOSE 7092

ENTRYPOINT ["/app/homeguard"]
CMD ["-http", ":7092"]

