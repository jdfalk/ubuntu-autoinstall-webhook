// filepath: /web/angular/autoinstall-web/src/app/components/logs/logs.component.ts
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { GrpcService } from '../../services/grpc.service';
import { TelemetryService } from '../../services/telemetry.service';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';

interface LogEntry {
    id: string;
    timestamp: Date;
    level: string;
    source: string;
    message: string;
    clientId?: string;
    metadata?: any;
}

@Component({
    selector: 'app-logs',
    templateUrl: './logs.component.html',
    styleUrls: ['./logs.component.scss']
})
export class LogsComponent implements OnInit {
    @ViewChild(MatPaginator) paginator!: MatPaginator;
    @ViewChild(MatSort) sort!: MatSort;

    displayedColumns: string[] = ['timestamp', 'level', 'source', 'message', 'actions'];
    dataSource = new MatTableDataSource<LogEntry>();
    logLevels: string[] = ['DEBUG', 'INFO', 'WARNING', 'ERROR', 'CRITICAL'];
    selectedLevels = new FormControl(this.logLevels);

    searchTerm = new FormControl('');
    sourceFilter = new FormControl('');
    startDate = new FormControl(null);
    endDate = new FormControl(null);

    isLoading = false;
    error = '';

    constructor(
        private grpcService: GrpcService,
        private telemetryService: TelemetryService
    ) { }

    ngOnInit(): void {
        this.setupFilters();
        this.loadLogs();
    }

    setupFilters(): void {
        // React to search term changes
        this.searchTerm.valueChanges.pipe(
            debounceTime(400),
            distinctUntilChanged()
        ).subscribe(() => this.applyFilters());

        // React to level selection changes
        this.selectedLevels.valueChanges.subscribe(() => this.applyFilters());

        // React to source filter changes
        this.sourceFilter.valueChanges.pipe(
            debounceTime(400),
            distinctUntilChanged()
        ).subscribe(() => this.applyFilters());

        // React to date changes
        this.startDate.valueChanges.subscribe(() => this.applyFilters());
        this.endDate.valueChanges.subscribe(() => this.applyFilters());
    }

    loadLogs(): void {
        this.isLoading = true;
        this.error = '';

        // Create mock data for demonstration
        setTimeout(() => {
            const mockLogs: LogEntry[] = [
                {
                    id: '1',
                    timestamp: new Date(Date.now() - 5 * 60000),
                    level: 'INFO',
                    source: 'api',
                    message: 'New client connected: server-web-01'
                },
                {
                    id: '2',
                    timestamp: new Date(Date.now() - 15 * 60000),
                    level: 'DEBUG',
                    source: 'installer',
                    message: 'Processing installation request for server-db-03'
                },
                {
                    id: '3',
                    timestamp: new Date(Date.now() - 20 * 60000),
                    level: 'WARNING',
                    source: 'auth',
                    message: 'Failed login attempt for user admin'
                },
                {
                    id: '4',
                    timestamp: new Date(Date.now() - 30 * 60000),
                    level: 'ERROR',
                    source: 'ipxe',
                    message: 'Failed to generate boot script: Template not found'
                },
                {
                    id: '5',
                    timestamp: new Date(Date.now() - 45 * 60000),
                    level: 'INFO',
                    source: 'installer',
                    message: 'Installation completed for server-api-02'
                },
                {
                    id: '6',
                    timestamp: new Date(Date.now() - 60 * 60000),
                    level: 'DEBUG',
                    source: 'config',
                    message: 'Configuration updated by user admin'
                },
                {
                    id: '7',
                    timestamp: new Date(Date.now() - 120 * 60000),
                    level: 'CRITICAL',
                    source: 'database',
                    message: 'Connection to database lost. Retrying...'
                }
            ];

            this.dataSource.data = mockLogs;
            this.dataSource.paginator = this.paginator;
            this.dataSource.sort = this.sort;
            this.isLoading = false;

            // Apply custom sort for date
            this.dataSource.sortingDataAccessor = (item, property) => {
                switch (property) {
                    case 'timestamp': return item.timestamp.getTime();
                    default: return (item as any)[property];
                }
            };

            // Set up filtering
            this.dataSource.filterPredicate = this.createFilter();
        }, 1000);

        // Real implementation would use gRPC:
        /*
        this.grpcService.callService('LogService', 'GetLogs', {
          levels: this.selectedLevels.value,
          search: this.searchTerm.value,
          source: this.sourceFilter.value,
          startTime: this.startDate.value ? new Date(this.startDate.value).toISOString() : null,
          endTime: this.endDate.value ? new Date(this.endDate.value).toISOString() : null,
          page: 0,
          pageSize: 100
        }).subscribe({
          next: (response) => {
            this.dataSource.data = response.logs.map((log: any) => ({
              ...log,
              timestamp: new Date(log.timestamp)
            }));
            this.dataSource.paginator = this.paginator;
            this.dataSource.sort = this.sort;
            this.isLoading = false;
          },
          error: (err) => {
            this.error = `Failed to load logs: ${err.message}`;
            this.isLoading = false;
          }
        });
        */
    }

    createFilter(): (data: LogEntry, filter: string) => boolean {
        return (data: LogEntry, filter: string): boolean => {
            // Check if the log level is in selected levels
            if (this.selectedLevels.value && !this.selectedLevels.value.includes(data.level)) {
                return false;
            }

            // Check source filter
            if (this.sourceFilter.value &&
                !data.source.toLowerCase().includes(this.sourceFilter.value.toLowerCase())) {
                return false;
            }

            // Check date range
            if (this.startDate.value && data.timestamp < new Date(this.startDate.value)) {
                return false;
            }

            if (this.endDate.value && data.timestamp > new Date(this.endDate.value)) {
                return false;
            }

            // Check search term against message
            if (this.searchTerm.value &&
                !data.message.toLowerCase().includes(this.searchTerm.value.toLowerCase())) {
                return false;
            }

            return true;
        };
    }

    applyFilters(): void {
        this.dataSource.filter = Math.random().toString(); // Trigger filter
    }

    resetFilters(): void {
        this.searchTerm.setValue('');
        this.selectedLevels.setValue(this.logLevels);
        this.sourceFilter.setValue('');
        this.startDate.setValue(null);
        this.endDate.setValue(null);
        this.applyFilters();
    }

    viewLogDetails(log: LogEntry): void {
        console.log('View log details:', log);
        // Open dialog with log details
    }

    downloadLogs(): void {
        const filteredData = this.dataSource.filteredData;

        if (filteredData.length === 0) {
            return;
        }

        // Prepare CSV content
        const headers = ['Timestamp', 'Level', 'Source', 'Message'];
        const rows = filteredData.map(log => [
            log.timestamp.toISOString(),
            log.level,
            log.source,
            log.message
        ]);

        const csvContent = [
            headers.join(','),
            ...rows.map(row => row.join(','))
        ].join('\n');

        // Create download link
        const blob = new Blob([csvContent], { type: 'text/csv' });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = url;
        a.download = `logs-${new Date().toISOString().substring(0, 10)}.csv`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    }
}
