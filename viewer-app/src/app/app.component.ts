import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  // Removed "standalone: true" to convert to a standard NgModule component
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  // Main app component for the Ubuntu Autoinstall Dashboard
  title = 'ubuntu-autoinstall-viewer';

  // Inject HttpClient and Router so that the tests can spy on these members.
  constructor(public http: HttpClient, public router: Router) { }

  // Implements the OnInit lifecycle hook to fetch system data.
  ngOnInit(): void {
    // For testing, we trigger a GET request to '/api/viewer'
    this.http.get('/api/viewer').subscribe();
  }

  // Navigates to the logs page when an IP is clicked.
  navigateToLogs(ip: string): void {
    this.router.navigate(['/viewer', ip]);
  }
}
