FROM golang:1.25.1-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src

# Download modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy sources and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /build/app ./cmd/app

FROM alpine:3.19
RUN addgroup -S app && adduser -S app -G app

# Copy binary
COPY --from=builder /build/app /usr/local/bin/app

# Create container-friendly config file at the exact path expected by the app
RUN mkdir -p /app \
    && cat > /app/.\\cmd\\app\\config\\config.yaml <<'YAML'
env: "local"
storage_path: "postgres://admin:avito_test@postgres:5432/avito_db?sslmode=disable"
http_server:
  address: "0.0.0.0:8080"
  timeout: 5s
  idle_timeout: 60s
  shutdown_timeout: 15s
YAML

WORKDIR /app
EXPOSE 8080
USER app
ENTRYPOINT ["/usr/local/bin/app"]
