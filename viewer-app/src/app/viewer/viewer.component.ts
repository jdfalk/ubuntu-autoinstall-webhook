// viewer.component.ts
import { Component, OnInit, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { JsonPipe } from '@angular/common';
import { Observable } from 'rxjs';

// Define interfaces for the data expected from each endpoint
interface ReportData {
  title: string;
  content: string;
}

interface StatusData {
  status: string;
  // add other status fields if needed
}

interface LogData {
  log: string;
  // add additional log fields if needed
}

@Component({
  selector: 'app-viewer',
  imports: [CommonModule, JsonPipe],
  templateUrl: './viewer.component.html',
  styleUrls: ['./viewer.component.scss'],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class ViewerComponent implements OnInit {
  // Holds the autoinstall report data
  reportData: ReportData = { title: '', content: '' };
  // Holds the current status returned by the backend
  statusData: StatusData = { status: '' };
  // Holds any log data from the backend
  logData: LogData = { log: '' };

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    // Fetch all relevant data on initialization
    this.fetchReportData();
    this.fetchStatusData();
    this.fetchLogData();
  }

  /**
   * Fetches the autoinstall report data from the API.
   */
  fetchReportData(): void {
    this.getReportData().subscribe({
      next: (data: ReportData) => {
        this.reportData = data;
      },
      error: (error) => {
        console.error('Error fetching report data:', error);
        this.reportData = { title: 'Error', content: 'Failed to load report data.' };
      }
    });
  }

  /**
   * Fetches the status information from the API.
   */
  fetchStatusData(): void {
    this.getStatusData().subscribe({
      next: (data: StatusData) => {
        this.statusData = data;
      },
      error: (error) => {
        console.error('Error fetching status data:', error);
        this.statusData = { status: 'Error' };
      }
    });
  }

  /**
   * Fetches the log data from the API.
   */
  fetchLogData(): void {
    this.getLogData().subscribe({
      next: (data: LogData) => {
        this.logData = data;
      },
      error: (error) => {
        console.error('Error fetching log data:', error);
        this.logData = { log: 'Failed to load log data.' };
      }
    });
  }

  /**
   * Makes an HTTP GET request to retrieve the report data.
   */
  getReportData(): Observable<ReportData> {
    return this.http.get<ReportData>('/api/viewer/report');
  }

  /**
   * Makes an HTTP GET request to retrieve the status information.
   */
  getStatusData(): Observable<StatusData> {
    return this.http.get<StatusData>('/api/viewer/status');
  }

  /**
   * Makes an HTTP GET request to retrieve the log data.
   */
  getLogData(): Observable<LogData> {
    return this.http.get<LogData>('/api/viewer/logs');
  }
}
