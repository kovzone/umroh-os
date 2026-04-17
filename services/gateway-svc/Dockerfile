# This Dockerfile must be run with the project root as context.
# e.g.: docker build -f services/gateway-svc/Dockerfile .

# Stage 1: Build environment
FROM golang:1.25-alpine AS build-env

# Set working directory for the build
WORKDIR /app

# Copy all projects based on go.mod replace directives so that go build sees them
COPY services/gateway-svc ./api

# Run inside gateway-svc
WORKDIR /app/api

# Download and verify dependencies
RUN go mod tidy
RUN go mod download

# Build the Go application
# CGO_ENABLED=0: Pure Go (no C dependencies)
# GOOS=linux: Target OS
# Compile all Go files in cmd directory into a single binary named 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/*.go

# Stage 2: Production environment
# Using minimal alpine image for the final container
FROM alpine:latest

# Install additional packages needed for the application
RUN apk update && apk upgrade && \
     apk add --no-cache bash git openssh

# Set the working directory for the application
WORKDIR /app

# Copy only necessary artifacts from the build stage
COPY --from=build-env /app/api/main main

# Note: config.json is mounted via docker-compose volume
# This allows for different configs (local vs docker) without rebuilding

# gateway-svc is REST-only (no gRPC server). Edge port 4000.
EXPOSE 4000

# Run the executable with 'start' argument
ENTRYPOINT [ "./main", "start" ]
