import { NgModule, APP_INITIALIZER } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';

// Angular Material
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTableModule } from '@angular/material/table';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';

// Monaco Editor
import { MonacoEditorModule } from 'ngx-monaco-editor-v2';

// App Components
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { LogsComponent } from './components/logs/logs.component';
import { IdeComponent } from './components/ide/ide.component';
import { ConfigComponent } from './components/config/config.component';
import { AccessDeniedComponent } from './components/access-denied/access-denied.component';
import { SystemsComponent } from './components/systems/systems.component';
import { InstallStatusComponent } from './components/install-status/install-status.component';
import { ApprovalComponent } from './components/approval/approval.component';
import { MonacoEditorComponent } from './components/monaco-editor/monaco-editor.component';
import { LoginComponent } from './components/login/login.component';

// Services
import { AuthService } from './services/auth.service';
import { GrpcService } from './services/grpc.service';
import { TelemetryService } from './services/telemetry.service';

// Guards
import { RoleGuard } from './guards/role.guard';

// Initialize telemetry on app start
export function initializeTelemetry(telemetryService: TelemetryService) {
  return () => { /* Telemetry is initialized in the constructor */ };
}

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    LogsComponent,
    IdeComponent,
    ConfigComponent,
    AccessDeniedComponent,
    SystemsComponent,
    InstallStatusComponent,
    ApprovalComponent,
    MonacoEditorComponent,
    LoginComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    AppRoutingModule,
    // Angular Material
    MatButtonModule,
    MatCardModule,
    MatDialogModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatMenuModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    MatSidenavModule,
    MatSnackBarModule,
    MatTableModule,
    MatTabsModule,
    MatToolbarModule,
    MatTooltipModule,
    // Monaco Editor
    MonacoEditorModule.forRoot()
  ],
  providers: [
    AuthService,
    GrpcService,
    RoleGuard,
    TelemetryService,
    {
      provide: APP_INITIALIZER,
      useFactory: initializeTelemetry,
      deps: [TelemetryService],
      multi: true
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
