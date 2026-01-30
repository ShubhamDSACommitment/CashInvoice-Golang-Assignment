# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./cmd/server

# ---------- Runtime Stage ----------
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose API port
EXPOSE 8080

# Run the server
CMD ["./server"]
