# Build stage
FROM golang:1.22.7 as builder

WORKDIR /app
COPY . .

# Build only the UI binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ui ./internal/ui/main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ui .

EXPOSE 8081
ENTRYPOINT ["./ui"]
