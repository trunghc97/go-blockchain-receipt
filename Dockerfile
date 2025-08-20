# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod file
COPY go.mod ./

# Download dependencies and generate go.sum
RUN go mod download && \
    go get -u github.com/swaggo/swag/cmd/swag && \
    go get -u github.com/swaggo/echo-swagger && \
    go get -u github.com/KyleBanks/depth && \
    go get -u github.com/go-openapi/jsonpointer && \
    go get -u github.com/go-openapi/jsonreference && \
    go get -u github.com/go-openapi/spec && \
    go get -u github.com/go-openapi/swag && \
    go get github.com/skip2/go-qrcode && \
    go get go.mongodb.org/mongo-driver/bson && \
    go get go.mongodb.org/mongo-driver/mongo && \
    go get go.mongodb.org/mongo-driver/mongo/options && \
    go mod tidy

# Copy source code
COPY . .

# Generate Swagger docs
RUN swag init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/app .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app"]
