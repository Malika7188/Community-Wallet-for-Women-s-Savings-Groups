# Build stage
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Runtime stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    bash \
    && rm -rf /var/lib/apt/lists/*

# Install Soroban CLI
RUN curl -L https://github.com/stellar/stellar-cli/releases/download/v22.0.1/stellar-cli-22.0.1-x86_64-unknown-linux-gnu.tar.gz -o stellar-cli.tar.gz && \
    tar -xzf stellar-cli.tar.gz && \
    mv stellar /usr/local/bin/stellar && \
    ln -s /usr/local/bin/stellar /usr/local/bin/soroban && \
    rm stellar-cli.tar.gz && \
    chmod +x /usr/local/bin/stellar

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy the Soroban contract WASM file
COPY --from=builder /app/chama_savings.wasm ./

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
