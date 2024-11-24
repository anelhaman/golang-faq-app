# Base image for building the Go app
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy Go dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the Go application source
COPY . .

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o main .

# Final image
FROM alpine:latest

# Install NGINX
RUN apk add --no-cache nginx ca-certificates && \
    mkdir -p /run/nginx

# Set the working directory
WORKDIR /app

# Copy the Go app binary and its dependencies
COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yaml config/
COPY --from=builder /app/data data/
COPY --from=builder /go/pkg/mod/github.com/narongdejsrn/go-thaiwordcut@v0.0.0-20190610123805-0a152d1829c4/dict/lexitron.txt /go/pkg/mod/github.com/narongdejsrn/go-thaiwordcut@v0.0.0-20190610123805-0a152d1829c4/dict/lexitron.txt

# Copy the NGINX configuration
COPY nginx.conf /etc/nginx/nginx.conf

ENV GIN_MODE=release

# Expose ports
EXPOSE 80 8080

# Command to start both services (Go app and NGINX)
CMD ["sh", "-c", "nginx && ./main"]
