import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-config-editor',
  templateUrl: './config-editor.component.html',
  styleUrls: ['./config-editor.component.css']
})
export class ConfigEditorComponent implements OnInit {
  // The config type can be 'ipxe' or 'cloudinit'
  configType: 'ipxe' | 'cloudinit' = 'ipxe';
  configContent: string = '';
  message: string = '';

  // Options for the Monaco Editor
  editorOptions = { theme: 'vs-dark', language: 'yaml', automaticLayout: true };

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    // Initially load ipxe config
    this.loadConfig();
  }

  // Load config based on selected type
  loadConfig(): void {
    const endpoint = this.configType === 'ipxe' ? '/api/ipxe-configs' : '/api/cloud-init-configs';
    this.http.get<any>(endpoint).subscribe(data => {
      // For simplicity, assume the config is in a field called "content"
      this.configContent = data.length > 0 ? data[0].content : '';
    }, err => {
      console.error('Error loading config', err);
      this.message = 'Error loading config';
    });
  }

  // Save updated config
  saveConfig(): void {
    const endpoint = this.configType === 'ipxe' ? '/api/ipxe-configs' : '/api/cloud-init-configs';
    // Here we assume a POST endpoint that accepts JSON with { content: string }
    this.http.post(endpoint, { content: this.configContent }).subscribe(response => {
      this.message = 'Config saved successfully!';
    }, err => {
      console.error('Error saving config', err);
      this.message = 'Error saving config';
    });
  }
}
