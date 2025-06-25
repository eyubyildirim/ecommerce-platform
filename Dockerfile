# --- Stage 1: Development Builder with Live-Reloading ---
# This stage contains the Go toolchain and 'air' for live reloading.
# We will use this stage for local development in Docker Compose.
FROM golang:1.24-alpine AS development

# Set the working directory inside the container
WORKDIR /app

# Install 'air' for live-reloading Go applications
RUN go install github.com/air-verse/air@latest

# Copy Go module files and download dependencies.
# This leverages Docker's layer caching. Dependencies are only re-downloaded
# if go.mod or go.sum change.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container
COPY . .

# The 'air' command will be the entry point for the development container.
# It will read the .air.toml file to know how to build and run the service.
CMD ["air"]


# --- Stage 2: Production Builder ---
# This stage builds the lean, final production binary.
FROM golang:1.24-alpine AS builder

# This ARG will be passed from docker-compose.yml to tell the builder
# WHICH service to build (e.g., "order_service").
ARG SERVICE_NAME

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the specified service. The output will be a single, static binary.
RUN CGO_ENABLED=0 go build -o /bin/${SERVICE_NAME} ./cmd/${SERVICE_NAME}/main/main.go


# --- Stage 3: Final Production Image ---
# This stage creates the final, minimal image for deployment.
FROM alpine:latest

ARG SERVICE_NAME

# Copy only the compiled binary from the 'builder' stage.
COPY --from=builder /bin/${SERVICE_NAME} /${SERVICE_NAME}

# Set the entry point for the container to be our compiled Go binary.
ENTRYPOINT [ "/" ]
