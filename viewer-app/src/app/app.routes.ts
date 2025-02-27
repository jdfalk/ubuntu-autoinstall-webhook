// app.routes.ts
import { Routes } from '@angular/router';
import { ViewerComponent } from './viewer/viewer.component';
import { ViewerListComponent } from './viewer-list/viewer-list.component';
import { ConfigEditorComponent } from './config-editor/config-editor.component';

export const routes: Routes = [
  { path: '', redirectTo: '/viewer-app/list', pathMatch: 'full' },
  { path: 'viewer-app/list', component: ViewerListComponent },
  { path: 'viewer', component: ViewerComponent },
  { path: 'config-editor', component: ConfigEditorComponent }
];
