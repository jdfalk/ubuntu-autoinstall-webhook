import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-system-log',
  templateUrl: './system-log.component.html',
  styleUrls: ['./system-log.component.css']
})
export class SystemLogsComponent implements OnInit {
  systemLogs: any[] = [];

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    // Fetch system logs from the database
    this.http.get<any[]>('/api/system-logs').subscribe(data => {
      this.systemLogs = data;
    }, err => {
      console.error('Error fetching system logs from DB', err);
    });
  }
}
