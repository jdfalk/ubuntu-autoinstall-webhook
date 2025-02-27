/**
 * AppModule
 *
 * This is the root Angular module for the application.
 * It imports essential modules such as BrowserModule and BrowserAnimationsModule,
 * sets up routing via AppRoutingModule, and declares application components.
 * Additionally, it configures Angular Material modules and the ngx-monaco-editor-v2.
 *
 * Angular Version: 19
 * Monaco Editor: ngx-monaco-editor-v2
 */

import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
// BrowserAnimationsModule is required for Angular Material animations and theming.
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { CommonModule } from '@angular/common';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ViewerComponent } from './viewer/viewer.component';
import { ConfigEditorComponent } from './config-editor/config-editor.component';
import { FormsModule } from '@angular/forms';
// Import the ngx-monaco-editor-v2 module and initialize it.
import { MonacoEditorModule } from 'ngx-monaco-editor-v2';

// Angular Material Modules
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatButtonModule } from '@angular/material/button';
import { MatListModule } from '@angular/material/list';

@NgModule({
    // Declare the components used in the application.
    declarations: [
        AppComponent,
        ViewerComponent,
        ConfigEditorComponent
    ],
    // Import necessary modules for the application.
    imports: [
        BrowserModule,
        BrowserAnimationsModule, // Required for Angular Material animations.
        CommonModule,
        AppRoutingModule,
        FormsModule,
        // Initialize ngx-monaco-editor-v2 with default configuration.
        MonacoEditorModule.forRoot(),
        // Angular Material modules for UI components.
        MatToolbarModule,
        MatSidenavModule,
        MatButtonModule,
        MatListModule
    ],
    providers: [],
    bootstrap: [AppComponent],
    // Allow the use of custom elements within the application.
    schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class AppModule { }
