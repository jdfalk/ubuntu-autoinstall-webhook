export const environment = {
    production: true,
    grpcEndpoint: '/grpc',
    apiEndpoint: '/api',
    useMutualTLS: true,     // Enable mutual TLS in production
    clientCert: '',         // This would be populated at deploy time or through a separate config
    telemetryEnabled: true,
    telemetryEndpoint: '/telemetry'  // In production, this might be a relative path
};
