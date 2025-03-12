// filepath: /web/angular/autoinstall-web/src/app/services/telemetry.service.ts
import { Injectable } from '@angular/core';
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
import { XMLHttpRequestInstrumentation } from '@opentelemetry/instrumentation-xml-http-request';
import { environment } from '../../environments/environment';

@Injectable({
    providedIn: 'root'
})
export class TelemetryService {
    private tracerProvider!: WebTracerProvider;

    constructor() {
        if (environment.telemetryEnabled) {
            this.initializeTracing();
        }
    }

    private initializeTracing(): void {
        const resource = Resource.default().merge(
            new Resource({
                [SemanticResourceAttributes.SERVICE_NAME]: 'autoinstall-web',
                [SemanticResourceAttributes.SERVICE_VERSION]: '1.0.0',
                [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]: environment.production ? 'production' : 'development'
            })
        );

        const exporter = new OTLPTraceExporter({
            url: environment.telemetryEndpoint + '/v1/traces'
        });

        this.tracerProvider = new WebTracerProvider({ resource });
        this.tracerProvider.addSpanProcessor(new BatchSpanProcessor(exporter));

        // Register instrumentations for automatic tracing
        registerInstrumentations({
            instrumentations: [
                new DocumentLoadInstrumentation(),
                new FetchInstrumentation({
                    propagateTraceHeaderCorsUrls: [/.*/]
                }),
                new XMLHttpRequestInstrumentation({
                    propagateTraceHeaderCorsUrls: [/.*/]
                })
            ],
            tracerProvider: this.tracerProvider
        });

        // Register the provider with a Zone context manager
        this.tracerProvider.register({
            contextManager: new ZoneContextManager()
        });
    }

    // Method to create custom spans
    createSpan(name: string, fn: () => any, attributes: Record<string, any> = {}): any {
        if (!environment.telemetryEnabled) {
            return fn();
        }

        const tracer = this.tracerProvider.getTracer('custom-tracer');
        const span = tracer.startSpan(name, { attributes });
        try {
            const result = fn();
            span.end();
            return result;
        } catch (error) {
            span.recordException(error as Error);
            span.setStatus({ code: 2, message: (error as Error).message });
            span.end();
            throw error;
        }
    }

    // Async version of createSpan
    async createAsyncSpan<T>(name: string, fn: () => Promise<T>, attributes: Record<string, any> = {}): Promise<T> {
        if (!environment.telemetryEnabled) {
            return fn();
        }

        const tracer = this.tracerProvider.getTracer('custom-tracer');
        const span = tracer.startSpan(name, { attributes });
        try {
            const result = await fn();
            span.end();
            return result;
        } catch (error) {
            span.recordException(error as Error);
            span.setStatus({ code: 2, message: (error as Error).message });
            span.end();
            throw error;
        }
    }
}
