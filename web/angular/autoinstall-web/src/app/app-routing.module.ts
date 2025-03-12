import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from './guards/role.guard';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { LogsComponent } from './components/logs/logs.component';
import { IdeComponent } from './components/ide/ide.component';
import { ConfigComponent } from './components/config/config.component';
import { AccessDeniedComponent } from './components/access-denied/access-denied.component';
import { SystemsComponent } from './components/systems/systems.component';
import { InstallStatusComponent } from './components/install-status/install-status.component';
import { ApprovalComponent } from './components/approval/approval.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'dashboard',
    pathMatch: 'full'
  },
  {
    path: 'dashboard',
    component: DashboardComponent
  },
  {
    path: 'logs',
    component: LogsComponent,
    canActivate: [RoleGuard],
    data: { requiredRole: 'logging' }
  },
  {
    path: 'ide',
    component: IdeComponent,
    canActivate: [RoleGuard],
    data: { requiredRole: 'ide' }
  },
  {
    path: 'config',
    component: ConfigComponent,
    canActivate: [RoleGuard],
    data: { requiredRole: 'configuration' }
  },
  {
    path: 'systems',
    component: SystemsComponent,
    canActivate: [RoleGuard],
    data: { requiredRole: 'admin' }
  },
  {
    path: 'install-status',
    component: InstallStatusComponent,
    canActivate: [RoleGuard],
    data: { requiredRole: 'admin' }
  },
  {
    path: 'approval',
    component: ApprovalComponent,
    canActivate: [RoleGuard],
    data: { requiredRole: 'admin' }
  },
  {
    path: 'access-denied',
    component: AccessDeniedComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
