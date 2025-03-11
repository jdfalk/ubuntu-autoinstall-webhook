import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-config-editor',
  templateUrl: './config-editor.component.html',
  styleUrls: ['./config-editor.component.scss']
})
export class ConfigEditorComponent implements OnInit {
  configContent: string = '';
  message: string = '';

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    this.loadConfig();
  }

  loadConfig(): void {
    // Replace with your actual backend endpoint
    this.http.get('http://localhost:8080/v1/config', { responseType: 'text' })
      .subscribe({
        next: (data) => this.configContent = data,
        error: (err) => {
          console.error('Error loading config:', err);
          this.message = 'Error loading configuration';
        }
      });
  }

  saveConfig(): void {
    // Replace with your actual backend endpoint
    this.http.post('http://localhost:8080/v1/config', { config: this.configContent })
      .subscribe({
        next: () => this.message = 'Configuration saved successfully!',
        error: (err) => {
          console.error('Error saving config:', err);
          this.message = 'Error saving configuration';
        }
      });
  }
}
