#######################################
# 1. Builder stage: compile the Go binary
#######################################
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod / go.sum first, so Docker can cache module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build a statically‐linked binary called "auth-service"
# Use CGO_ENABLED=0 to force a pure Go binary.
RUN CGO_ENABLED=0 go build -o /auth-service

#######################################
# 2. Final stage: minimal runtime image
#######################################
FROM alpine:latest


WORKDIR /root

# Copy the compiled binary from the builder stage
COPY --from=builder /auth-service .
# Expose port 80 (or whatever your app uses)
EXPOSE 8090

# By default, run the binary. It will read PORT from environment or .env.
ENTRYPOINT ["./auth-service"]
