import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-ipxe-config',
  templateUrl: './ipxe-config.component.html',
  styleUrls: ['./ipxe-config.component.css']
})
export class IpxeConfigComponent implements OnInit {
  ipxeConfigs: any[] = [];
  historicalConfigs: any[] = [];

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    // Fetch current ipxe configs from the database
    this.http.get<any[]>('/api/ipxe-configs').subscribe(data => {
      this.ipxeConfigs = data;
    }, err => {
      console.error('Error fetching ipxe configs', err);
    });

    // Fetch historical ipxe configs from the database
    this.http.get<any[]>('/api/historical-ipxe-configs').subscribe(data => {
      this.historicalConfigs = data;
    }, err => {
      console.error('Error fetching historical ipxe configs', err);
    });
  }
}
