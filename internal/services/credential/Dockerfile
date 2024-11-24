# Build stage
FROM golang:1.22.7 as builder

WORKDIR /app
COPY . .

# Build only the credential service binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o credential-service ./internal/services/credential/main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/credential-service .
COPY ./migrations /app/migrations

EXPOSE 8080
ENTRYPOINT ["./credential-service"]
