FROM alpine:latest

ARG PB_VERSION=0.38.1

# Install necessary certificates and tools
RUN apk add --no-cache ca-certificates unzip wget

# Download and unzip PocketBase
RUN wget https://github.com/pocketbase/pocketbase/releases/download/v${PB_VERSION}/pocketbase_${PB_VERSION}_linux_amd64.zip \
    && unzip pocketbase_${PB_VERSION}_linux_amd64.zip \
    && chmod +x /pocketbase

# Expose the default PocketBase port
EXPOSE 8090

# Start PocketBase, binding to all network interfaces
CMD ["/pocketbase", "serve", "--http=0.0.0.0:8090"]