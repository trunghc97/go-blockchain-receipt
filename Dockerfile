# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Copy entire project
COPY . .

# Install swag and generate docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init

# Download dependencies and build
RUN go mod download && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and swagger docs from builder
COPY --from=builder /app/app .
COPY --from=builder /app/api/docs ./api/docs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app"]
