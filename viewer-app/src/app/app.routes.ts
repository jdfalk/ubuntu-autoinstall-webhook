import { Routes } from '@angular/router';
import { AppComponent } from './app.component';

export const routes: Routes = [
  { path: '', redirectTo: '/viewer', pathMatch: 'full' },
  { path: 'viewer', component: AppComponent },
  { path: 'viewer/:ip', component: AppComponent }
];
