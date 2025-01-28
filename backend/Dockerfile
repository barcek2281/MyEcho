# Stage 1: Build
FROM golang:1.23.3 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN make

# Stage 2: Final image
FROM alpine:latest

# Install required packages for PostgreSQL
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary from the builder stage
COPY --from=builder /app/MyEcho .

# configs
COPY ./configs/apiserver.toml /root/configs/apiserver.toml
COPY ./log/info.log /root/log/info.log
COPY ./templates /root/templates
COPY ./static /root/static
COPY ./migrations/1.up.sql /root/migrations/1.up.sql
COPY ./migrations/1.down.sql /root/migrations/1.down.sql

# Expose the application port
EXPOSE 8080
RUN go test -v ./...

# Run the binary
CMD ["./MyEcho"]


