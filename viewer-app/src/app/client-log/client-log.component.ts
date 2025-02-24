import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-client-log',
  templateUrl: './client-log.component.html',
  styleUrls: ['./client-log.component.css']
})
export class ClientLogsComponent implements OnInit {
  clientLogs: any[] = [];

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    // Fetch client logs from the database
    this.http.get<any[]>('/api/client-logs').subscribe(data => {
      this.clientLogs = data;
    }, err => {
      console.error('Error fetching client logs', err);
    });
  }
}
