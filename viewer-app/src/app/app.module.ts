import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';               // Import FormsModule for ngModel
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http'; // New HTTP provider API
import { MonacoEditorModule } from 'ngx-monaco-editor-v2';    // Import MonacoEditorModule from the new package

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { SystemLogsComponent } from './system-log/system-log.component';
import { ClientLogsComponent } from './client-log/client-log.component';
import { IpxeConfigComponent } from './ipxe-config/ipxe-config.component';
import { CloudInitConfigComponent } from './cloud-init-config/cloud-init-config.component';
import { ConfigEditorComponent } from './config-editor/config-editor.component';

@NgModule({
  declarations: [
    AppComponent,
    SystemLogsComponent,
    ClientLogsComponent,
    IpxeConfigComponent,
    CloudInitConfigComponent,
    ConfigEditorComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,                        // Add FormsModule here for template-driven forms
    MonacoEditorModule.forRoot()        // Initialize the Monaco Editor module using the correct export
  ],
  providers: [
    provideHttpClient(withInterceptorsFromDi())
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
