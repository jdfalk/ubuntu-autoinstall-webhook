import { Component, OnInit } from '@angular/core';
import { InstallService, StatusRequest, StatusResponse } from '../services/install.service';
import { CommonModule, JsonPipe } from '@angular/common';

@Component({
  selector: 'app-status-update',
  templateUrl: './status-update.component.html',
  imports: [JsonPipe, CommonModule],
  styleUrls: ['./status-update.component.scss']
})
export class StatusUpdateComponent implements OnInit {
  response: StatusResponse | null = null;

  constructor(private installService: InstallService) { }

  ngOnInit(): void { }

  sendStatus(): void {
    const status: StatusRequest = {
      hostname: 'test-host',
      ip_address: '192.168.1.100',
      progress: 50,
      message: 'Installation is halfway complete'
    };
    this.installService.reportStatus(status).subscribe({
      next: (res) => this.response = res,
      error: (err) => console.error('Error reporting status:', err)
    });
  }
}
