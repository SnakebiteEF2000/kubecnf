# Build stage
FROM golang:1.24-alpine AS builder

# Set build arguments
ARG VERSION=dev
ARG BUILDTIME

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates tzdata

# Create appuser for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILDTIME}" \
    -a -installsuffix cgo \
    -o kubecnf .

# Final stage
FROM scratch

# Copy ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy user from builder
COPY --from=builder /etc/passwd /etc/passwd

# Copy binary from builder
COPY --from=builder /build/kubecnf /kubecnf

# Use appuser
USER appuser

# Set entrypoint
ENTRYPOINT ["/kubecnf"] 