# --- Stage 1: Create Minimal Final Image ---
    FROM alpine:latest

    # Set build args for multi-arch
    ARG TARGETARCH

    # Set environment variables
    ENV WEBHOOK_PORT=25000 \
        CONFIG_FILE=/config/config.yaml

    WORKDIR /app

    # Copy the pre-built binary from the GitHub workflow
    COPY webhook-linux-${TARGETARCH} /usr/local/bin/webhook

    # Ensure correct permissions
    RUN chmod +x /usr/local/bin/webhook

    # Expose webhook port
    EXPOSE $WEBHOOK_PORT

    # Start the application
    ENTRYPOINT ["/usr/local/bin/webhook"]
