import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { GrpcService } from '../../services/grpc.service';
import { TelemetryService } from '../../services/telemetry.service';
import { Subscription, timer } from 'rxjs';
import { switchMap } from 'rxjs/operators';

interface SystemStats {
  totalSystems: number;
  activeSystems: number;
  pendingApprovals: number;
  failedInstalls: number;
}

interface SystemEvent {
  id: string;
  timestamp: Date;
  hostname: string;
  eventType: string;
  message: string;
}

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit, OnDestroy {
  isAdmin = false;
  stats: SystemStats = {
    totalSystems: 0,
    activeSystems: 0,
    pendingApprovals: 0,
    failedInstalls: 0
  };
  recentEvents: SystemEvent[] = [];
  loading = true;
  error = '';
  private refreshSubscription?: Subscription;

  constructor(
    private authService: AuthService,
    private grpcService: GrpcService,
    private telemetryService: TelemetryService,
    private router: Router
  ) { }

  ngOnInit(): void {
    // Check if user has admin role
    this.isAdmin = this.authService.hasRole('admin');

    // Load dashboard data
    this.loadDashboardData();

    // Set up automatic refresh every 30 seconds
    this.refreshSubscription = timer(30000, 30000).subscribe(() => {
      this.loadDashboardData(false);
    });
  }

  ngOnDestroy(): void {
    if (this.refreshSubscription) {
      this.refreshSubscription.unsubscribe();
    }
  }

  loadDashboardData(showLoading = true): void {
    if (showLoading) {
      this.loading = true;
    }
    this.error = '';

    // Using the telemetry service to create a span for this operation
    this.telemetryService.createAsyncSpan('dashboard.loadData', async () => {
      try {
        // For now, we're using mock data
        // In a real implementation, this would call the grpcService

        // Simulate API call delay
        await new Promise(resolve => setTimeout(resolve, 1000));

        // Mock stats data
        this.stats = {
          totalSystems: 42,
          activeSystems: 36,
          pendingApprovals: 3,
          failedInstalls: 2
        };

        // Mock events data
        this.recentEvents = [
          {
            id: '1',
            timestamp: new Date(Date.now() - 5 * 60000),
            hostname: 'server-web-01',
            eventType: 'InstallSuccess',
            message: 'Ubuntu 22.04 LTS installation completed successfully'
          },
          {
            id: '2',
            timestamp: new Date(Date.now() - 15 * 60000),
            hostname: 'server-db-03',
            eventType: 'InstallStarted',
            message: 'Ubuntu 22.04 LTS installation initiated'
          },
          {
            id: '3',
            timestamp: new Date(Date.now() - 20 * 60000),
            hostname: 'server-api-02',
            eventType: 'ApprovalRequired',
            message: 'New system requires approval for installation'
          }
        ];

        this.loading = false;

        // In real implementation:
        // this.grpcService.callService('DashboardService', 'GetStats', {})
        //   .pipe(
        //     switchMap(statsResponse => {
        //       this.stats = statsResponse;
        //       return this.grpcService.callService('DashboardService', 'GetRecentEvents', { limit: 5 });
        //     })
        //   )
        //   .subscribe({
        //     next: (eventsResponse) => {
        //       this.recentEvents = eventsResponse.events;
        //       this.loading = false;
        //     },
        //     error: (err) => {
        //       this.loading = false;
        //       this.error = `Error loading dashboard data: ${err.message}`;
        //     }
        //   });
      } catch (err) {
        this.loading = false;
        this.error = `Error loading dashboard data: ${(err as Error).message}`;
      }
    }, { component: 'dashboard' });
  }

  navigateTo(route: string): void {
    this.router.navigate([route]);
  }

  formatEventTime(date: Date): string {
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);

    if (diffMins < 1) {
      return 'Just now';
    } else if (diffMins < 60) {
      return `${diffMins} minute${diffMins !== 1 ? 's' : ''} ago`;
    } else {
      const diffHours = Math.floor(diffMins / 60);
      return `${diffHours} hour${diffHours !== 1 ? 's' : ''} ago`;
    }
  }
}
