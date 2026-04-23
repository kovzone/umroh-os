# This Dockerfile must be run with the project root as context.
# e.g.: docker build -f services/logistics-svc/Dockerfile .

# Stage 1: Build environment
FROM golang:1.25-alpine AS build-env

# Set working directory for the build
WORKDIR /app

# Copy all projects based on go.mod replace directives so that go build sees them
COPY services/logistics-svc ./api

# Run inside logistics-svc
WORKDIR /app/api

# Download and verify dependencies
RUN go mod tidy
RUN go mod download

# Build the Go application
# CGO_ENABLED=0: Pure Go (no C dependencies)
# GOOS=linux: Target OS
# Compile all Go files in cmd directory into a single binary named 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Download grpc_health_probe in build stage (golang:alpine already has wget)
ARG GRPC_HEALTH_PROBE_VERSION=v0.4.28
RUN wget -qO /bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# Stage 2: Production environment
# Using minimal alpine image for the final container
FROM alpine:latest

# Install additional packages needed for the application
RUN apk update && apk upgrade && \
     apk add --no-cache bash git openssh wget

# Set the working directory for the application
WORKDIR /app

# Copy only necessary artifacts from the build stage
COPY --from=build-env /app/api/main main
# Copy grpc_health_probe for docker-compose gRPC healthcheck (ADR 0009 / BL-MON-001).
COPY --from=build-env /bin/grpc_health_probe /bin/grpc_health_probe

# Note: config.json is mounted via docker-compose volume
# This allows for different configs (local vs docker) without rebuilding

# Expose gRPC (50058). Per ADR 0009 logistics-svc is gRPC-only — REST was removed
# in BL-REFACTOR-005 / S1-E-13.
EXPOSE 50058

# Run the executable with 'start' argument
ENTRYPOINT [ "./main", "start" ]
