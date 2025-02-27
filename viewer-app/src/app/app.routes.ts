// app.routes.ts
import { Routes } from '@angular/router';
import { ViewerComponent } from './viewer/viewer.component';
import { ConfigEditorComponent } from './config-editor/config-editor.component';

export const routes: Routes = [
  { path: '', redirectTo: '/viewer', pathMatch: 'full' },
  { path: 'viewer', component: ViewerComponent },
  { path: 'config-editor', component: ConfigEditorComponent }
];
