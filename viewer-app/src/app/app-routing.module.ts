import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SystemLogsComponent } from './system-log/system-log.component';
import { ClientLogsComponent } from './client-log/client-log.component';
import { IpxeConfigComponent } from './ipxe-config/ipxe-config.component';
import { CloudInitConfigComponent } from './cloud-init-config/cloud-init-config.component';
import { ConfigEditorComponent } from './config-editor/config-editor.component';

const routes: Routes = [
  { path: '', redirectTo: '/system-logs', pathMatch: 'full' },
  { path: 'system-logs', component: SystemLogsComponent },
  { path: 'client-logs', component: ClientLogsComponent },
  { path: 'ipxe-configs', component: IpxeConfigComponent },
  { path: 'cloud-init-configs', component: CloudInitConfigComponent },
  { path: 'config-editor', component: ConfigEditorComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
