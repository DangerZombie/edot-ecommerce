# Build Stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy entire source code
COPY . .

# Set working directory to service folder and build
WORKDIR /app/micro-services/user
RUN go build -o /bin/service .

# Runtime Stage
FROM ubuntu:22.04

# Copy binary and migrations
COPY --from=builder /bin/service /bin/service
COPY --from=builder /app/micro-services/user/migrations /usr/bin/migrations

# Add access
RUN chmod -R 644 /usr/bin/migrations

# Set working directory and entry point
WORKDIR /usr/bin
ENTRYPOINT ["/bin/service"]
