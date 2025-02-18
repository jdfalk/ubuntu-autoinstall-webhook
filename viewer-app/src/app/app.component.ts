import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-viewer',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  systems: any[] = [];
  logs: any[] = [];
  selectedIp: string = '';
  sortField: string = '';
  ascending: boolean = true;

  constructor(private http: HttpClient, private route: ActivatedRoute, private router: Router) {}

  ngOnInit() {
    this.http.get('/api/viewer').subscribe((data: any) => this.systems = data);
    this.route.params.subscribe(params => {
      if (params['ip']) {
        this.loadLogs(params['ip']);
      }
    });
  }

  navigateToLogs(ip: string) {
    this.router.navigate(['/viewer', ip]);
  }

  loadLogs(ip: string) {
    this.selectedIp = ip;
    this.http.get(`/api/viewer/${ip}`).subscribe((data: any) => this.logs = data);
  }

  sort(field: string) {
    if (this.sortField === field) {
      this.ascending = !this.ascending;
    } else {
      this.sortField = field;
      this.ascending = true;
    }
    this.logs.sort((a, b) => {
      return (this.ascending ? 1 : -1) * ((a[field] > b[field]) ? 1 : -1);
    });
  }
}
