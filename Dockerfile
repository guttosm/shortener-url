# Stage 1: Build the Go binary
FROM golang:1.24.2-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project
COPY . .

# Build the app from cmd/main.go
RUN go build -o url-shortener ./cmd/main.go

# Stage 2: Lightweight container
FROM alpine:3.21.3

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/url-shortener .

EXPOSE 8080

ENTRYPOINT ["./url-shortener"]