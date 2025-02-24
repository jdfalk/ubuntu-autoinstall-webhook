import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-cloud-init-config',
  templateUrl: './cloud-init-config.component.html',
  styleUrls: ['./cloud-init-config.component.css']
})
export class CloudInitConfigComponent implements OnInit {
  cloudInitConfigs: any[] = [];
  historicalConfigs: any[] = [];

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    // Fetch current cloud-init configs
    this.http.get<any[]>('/api/cloud-init-configs').subscribe(data => {
      this.cloudInitConfigs = data;
    }, err => {
      console.error('Error fetching cloud-init configs', err);
    });

    // Fetch historical cloud-init configs
    this.http.get<any[]>('/api/historical-cloud-init-configs').subscribe(data => {
      this.historicalConfigs = data;
    }, err => {
      console.error('Error fetching historical cloud-init configs', err);
    });
  }
}
