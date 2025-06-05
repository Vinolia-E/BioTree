# Build stage
FROM golang:1.23-alpine AS build

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o biotree

# Final stage
FROM alpine:latest

# Install required dependencies for document processing
RUN apk add --no-cache libreoffice

# Create directories for file uploads and data
RUN mkdir -p /app/files /app/data

# Set working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/biotree .

# Copy frontend assets
COPY --from=build /app/frontend ./frontend

# Create a non-root user
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./biotree"]