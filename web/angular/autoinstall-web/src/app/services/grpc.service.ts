import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
    providedIn: 'root'
})
export class GrpcService {
    private client: any;

    constructor() {
        // Initialize gRPC client with authentication
        this.initializeClient();
        // Access environment variables
        console.log('Using gRPC endpoint:', environment.grpcEndpoint);
        console.log('Mutual TLS enabled:', environment.useMutualTLS);
    }

    private initializeClient(): void {
        // Implementation will depend on your gRPC web client library
        // Example assuming grpc-web
        const { GrpcWebClientBase } = require('grpc-web');

        // Configure mutual TLS if available
        const useMTLS = environment.useMutualTLS;
        const options = useMTLS ? {
            withCredentials: true,
            metadata: {
                'X-Client-Certificate': environment.clientCert
            }
        } : {};

        this.client = new GrpcWebClientBase({
            url: environment.grpcEndpoint,
            ...options
        });
    }

    // Method to call specific gRPC services
    callService(service: string, method: string, request: any): Observable<any> {
        // Implementation depends on your specific protocol buffers
        return new Observable(observer => {
            // Example implementation
            const serviceClient = this.getServiceClient(service);
            serviceClient[method](request, (err: any, response: any) => {
                if (err) {
                    observer.error(err);
                } else {
                    observer.next(response);
                }
                observer.complete();
            });
        });
    }

    private getServiceClient(service: string): any {
        // Return the appropriate service client based on your proto definitions
        // This is a placeholder - real implementation would use your generated clients
        return null;
    }
}
