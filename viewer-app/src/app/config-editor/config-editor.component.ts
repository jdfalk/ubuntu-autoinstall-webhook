// config-editor.component.ts
import { Component, OnInit, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';


@Component({
  selector: 'app-config-editor',
  templateUrl: './config-editor.component.html',
  styleUrls: ['./config-editor.component.scss'],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class ConfigEditorComponent implements OnInit {
  editorOptions = {
    theme: 'vs-dark',
    language: 'json',
    automaticLayout: true
  };

  configContent: string = `{
  "setting1": "value1",
  "setting2": "value2"
}`;

  constructor() { }

  ngOnInit(): void {
    // Optionally load configuration from an API endpoint
  }

  onConfigChanged(newValue: any): void {
    // Cast the incoming value to string if needed
    const value = newValue as string;
    this.configContent = value;
    console.log('Configuration updated:', this.configContent);
  }
}
