import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

interface NodeInfo {
  id: string;
  hostname: string;
  mac: string;
  ip: string;
  ipxeDestination: string; // "install" or "os"
  cloudInitPhase: string; // "install" or "post-install"
}

@Component({
  selector: 'app-viewer-list',
  templateUrl: './viewer-list.component.html',
  styleUrls: ['./viewer-list.component.scss']
})
export class ViewerListComponent implements OnInit {
  nodes: NodeInfo[] = [
    // Sample data; in production you would fetch this from your API
    {
      id: 'node-1',
      hostname: 'node-1-host',
      mac: 'aa:bb:cc:dd:ee:ff',
      ip: '192.168.1.100',
      ipxeDestination: 'install',
      cloudInitPhase: 'install'
    },
    {
      id: 'node-2',
      hostname: 'node-2-host',
      mac: '11:22:33:44:55:66',
      ip: '192.168.1.101',
      ipxeDestination: 'os',
      cloudInitPhase: 'post-install'
    }
  ];

  constructor(private router: Router) {}

  ngOnInit(): void {}

  viewNode(node: NodeInfo): void {
    // Navigate to the viewer page with node id parameter
    this.router.navigate(['/viewer-app/viewer'], { queryParams: { node: node.id } });
  }
}
