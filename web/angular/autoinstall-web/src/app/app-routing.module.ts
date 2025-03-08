import { Routes } from '@angular/router';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ConfigEditorComponent } from './config-editor/config-editor.component';
import { StatusUpdateComponent } from './status-update/status-update.component';

export const appRoutes: Routes = [
    { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
    { path: 'dashboard', component: DashboardComponent },
    { path: 'config-editor', component: ConfigEditorComponent },
    { path: 'status-update', component: StatusUpdateComponent }
];

export const AppRoutingModule: Routes = [
    { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
    { path: 'dashboard', component: DashboardComponent },
    { path: 'config-editor', component: ConfigEditorComponent },
    { path: 'status-update', component: StatusUpdateComponent }
];
