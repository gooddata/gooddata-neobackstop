# Stage 0: Build the base go system with playwright drivers
FROM 020413372491.dkr.ecr.us-east-1.amazonaws.com/pullthrough/docker.io/library/golang:1.25.6-bookworm AS basesystem
WORKDIR /

# Install CA certs & dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    && update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Install Playwright drivers
# The CLI version MUST match the playwright-go version pinned in go.mod, otherwise
# the installed driver/browsers can be incompatible with the library and fail to launch.
RUN GOOS=linux GOARCH=${TARGETARCH} go run github.com/mxschmitt/playwright-go/cmd/playwright@v0.6201.1 install --with-deps chromium firefox

# Stage 1: Build the Go app
FROM 020413372491.dkr.ecr.us-east-1.amazonaws.com/pullthrough/docker.io/library/golang:1.25.6-bookworm AS builder
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Declare the TARGETARCH argument (automatically set by Docker BuildKit)
ARG TARGETARCH

# Build the Go application
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=${TARGETARCH} \
    go build -o bin .

# Stage 2: Run the app
FROM basesystem
# Copy the built Go app binary and the html_report_assets which are needed at runtime
RUN mkdir -p /usr/neobackstop/app
COPY --from=builder /app/bin /usr/neobackstop/app/bin
COPY --from=builder /app/html_report_assets /usr/neobackstop/app/html_report_assets

ENV ENVIRONMENT=PROD

WORKDIR /usr/neobackstop/app

# Use entrypoint to allow argument passing
ENTRYPOINT ["/usr/neobackstop/app/bin"]
# CMD will be default command if none provided
CMD ["test"]
