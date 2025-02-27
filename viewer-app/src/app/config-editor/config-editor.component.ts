// config-editor.component.ts
import { Component, OnInit, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';

@Component({
  selector: 'app-config-editor',
  templateUrl: './config-editor.component.html',
  styleUrls: ['./config-editor.component.scss'],
  // Keep the custom elements schema, and ensure Material modules are also imported in the NgModule.
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class ConfigEditorComponent implements OnInit {
  // Editor options for ngx-monaco-editor
  editorOptions = {
    theme: 'vs-dark',
    language: 'json',
    automaticLayout: true
  };

  // The JSON config text
  configContent: string = `{
  "setting1": "value1",
  "setting2": "value2"
}`;

  constructor() { }

  ngOnInit(): void {
    // Optionally load configuration from an API endpoint
    // E.g., fetch('your-api/config').then(...)
  }

  onConfigChanged(newValue: any): void {
    // Cast the incoming value to string if needed
    const value = newValue as string;
    this.configContent = value;
    console.log('Configuration updated:', this.configContent);
  }

  // Example button handler for "Save"
  onSave(): void {
    console.log('Save button clicked. Current config:', this.configContent);
    // TODO: Implement actual save logic here
  }

  // Example button handler for "Cancel"
  onCancel(): void {
    console.log('Cancel button clicked. Navigating away or discarding changes...');
    // TODO: Implement navigation or discard changes
  }
}
