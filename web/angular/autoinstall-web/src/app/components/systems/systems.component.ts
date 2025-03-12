// filepath: /Users/jdfalk/repos/github.com/jdfalk/ubuntu-autoinstall-webhook/web/angular/autoinstall-web/src/app/components/systems/systems.component.ts
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { GrpcService } from '../../services/grpc.service';

interface System {
    id: string;
    hostname: string;
    macAddress: string;
    ipAddress: string;
    status: string;
    lastSeen: Date;
}

@Component({
    selector: 'app-systems',
    templateUrl: './systems.component.html',
    styleUrls: ['./systems.component.scss']
})
export class SystemsComponent implements OnInit {
    systems: System[] = [];
    loading = true;
    error = '';

    constructor(
        private grpcService: GrpcService,
        private router: Router
    ) { }

    ngOnInit(): void {
        this.loadSystems();
    }

    loadSystems(): void {
        this.loading = true;
        this.error = '';

        this.grpcService.callService('ConfigurationService', 'ListSystems', {})
            .subscribe({
                next: (response) => {
                    this.systems = response.systems.map((sys: any) => ({
                        id: sys.id,
                        hostname: sys.hostname,
                        macAddress: sys.macAddress,
                        ipAddress: sys.ipAddress,
                        status: sys.status,
                        lastSeen: new Date(sys.lastSeen * 1000)
                    }));
                    this.loading = false;
                },
                error: (err) => {
                    this.loading = false;
                    this.error = `Error loading systems: ${err.message}`;
                    console.error('Error loading systems:', err);
                }
            });
    }

    viewLogs(system: System): void {
        this.router.navigate(['/logs'], { queryParams: { clientId: system.id } });
    }

    editConfig(system: System): void {
        this.router.navigate(['/ide'], { queryParams: { clientId: system.id, configType: 'ipxe' } });
    }

    editCloudInit(system: System): void {
        this.router.navigate(['/ide'], { queryParams: { clientId: system.id, configType: 'cloud-init' } });
    }
}
