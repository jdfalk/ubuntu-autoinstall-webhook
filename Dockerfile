# --- Stage 1: Minimal Image Using Official Go Alpine Base ---
    FROM golang:1.24-alpine AS final

    # Set build argument for multi-arch support
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
