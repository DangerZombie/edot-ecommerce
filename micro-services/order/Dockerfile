# Build Stage
FROM golang:1.22 AS builder

# Specify the working directory in the container
WORKDIR /app

# Copy the go.mod and go.sum files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Go to the appropriate service directory
WORKDIR /app/micro-services/order

# Build Go applications for order-service
RUN go build -o /bin/service .

# Runtime Stage
FROM ubuntu:22.04

# Copy the binary that has been built
COPY --from=builder /bin/service /bin/service

# Copy the migrations folder from the builder stage into the runtime image
COPY ./micro-services/order/migrations /app/micro-services/order/migrations

# Determine the entry point
ENTRYPOINT ["/bin/service"]
