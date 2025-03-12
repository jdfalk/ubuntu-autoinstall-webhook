// filepath: /web/angular/autoinstall-web/src/app/components/ide/ide.component.ts
import { Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { GrpcService } from '../../services/grpc.service';
import { TelemetryService } from '../../services/telemetry.service';
import { MonacoEditorComponent } from '../monaco-editor/monaco-editor.component';

@Component({
    selector: 'app-ide',
    templateUrl: './ide.component.html',
    styleUrls: ['./ide.component.scss']
})
export class IdeComponent implements OnInit {
    @ViewChild(MonacoEditorComponent) editor!: MonacoEditorComponent;

    clientId: string = '';
    configType: string = 'ipxe';
    configName: string = '';
    configContent: string = '';
    configLanguage: string = 'shell';
    isLoading: boolean = false;
    isSaving: boolean = false;
    error: string = '';
    availableConfigs: string[] = [];

    constructor(
        private route: ActivatedRoute,
        private grpcService: GrpcService,
        private snackBar: MatSnackBar,
        private telemetryService: TelemetryService
    ) { }

    ngOnInit(): void {
        this.route.queryParams.subscribe(params => {
            if (params['clientId']) {
                this.clientId = params['clientId'];
            }

            if (params['configType']) {
                this.configType = params['configType'];
                this.setLanguageByConfigType();
            }

            this.loadAvailableConfigs();
        });
    }

    loadAvailableConfigs(): void {
        this.isLoading = true;
        this.error = '';

        // Mock data for now
        setTimeout(() => {
            if (this.configType === 'ipxe') {
                this.availableConfigs = ['default.ipxe', 'server.ipxe', 'workstation.ipxe'];
            } else if (this.configType === 'cloud-init') {
                this.availableConfigs = ['user-data', 'meta-data', 'network-config'];
            } else {
                this.availableConfigs = ['config.yaml', 'templates.yaml'];
            }
            this.isLoading = false;

            // Load first config by default
            if (this.availableConfigs.length > 0) {
                this.loadConfig(this.availableConfigs[0]);
            }
        }, 500);

        // Real implementation would use gRPC:
        /*
        this.grpcService.callService('ConfigService', 'ListConfigs', {
          configType: this.configType,
          clientId: this.clientId || ''
        }).subscribe({
          next: (response) => {
            this.availableConfigs = response.configs || [];
            this.isLoading = false;

            // Load first config by default
            if (this.availableConfigs.length > 0) {
              this.loadConfig(this.availableConfigs[0]);
            }
          },
          error: (err) => {
            this.error = `Failed to load available configurations: ${err.message}`;
            this.isLoading = false;
          }
        });
        */
    }

    loadConfig(configName: string): void {
        this.configName = configName;
        this.isLoading = true;
        this.error = '';

        // Mock data for different config types
        setTimeout(() => {
            if (this.configType === 'ipxe') {
                this.configContent = `#!ipxe
# ${configName}
# This is a sample iPXE configuration

set base-url http://192.168.1.1/ubuntu
kernel \${base-url}/vmlinuz
initrd \${base-url}/initrd
imgargs vmlinuz initrd=initrd root=/dev/ram0 url=http://192.168.1.1/autoinstall.yaml
boot`;
            } else if (this.configType === 'cloud-init' && configName === 'user-data') {
                this.configContent = `#cloud-config
hostname: ubuntu-server
users:
  - name: ubuntu
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    ssh_authorized_keys:
      - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ...`;
            } else {
                this.configContent = `# ${configName} configuration
version: 1

settings:
  debug: true
  log_level: info`;
            }

            this.setLanguageByConfigType();
            this.isLoading = false;
        }, 800);

        // Real implementation would use gRPC:
        /*
        this.grpcService.callService('ConfigService', 'GetConfig', {
          configType: this.configType,
          configName: configName,
          clientId: this.clientId || ''
        }).subscribe({
          next: (response) => {
            this.configContent = response.content || '';
            this.setLanguageByConfigType();
            this.isLoading = false;
          },
          error: (err) => {
            this.error = `Failed to load configuration: ${err.message}`;
            this.isLoading = false;
          }
        });
        */
    }

    saveConfig(): void {
        if (!this.configName) {
            this.snackBar.open('Please select a configuration to save', 'Dismiss', { duration: 3000 });
            return;
        }

        this.isSaving = true;

        // Get content from Monaco editor
        const content = this.editor ? this.editor.getValue() : this.configContent;

        // Mock saving - just show success message
        setTimeout(() => {
            this.isSaving = false;
            this.snackBar.open('Configuration saved successfully', 'Dismiss', { duration: 3000 });
        }, 1000);

        // Real implementation would use gRPC:
        /*
        this.grpcService.callService('ConfigService', 'SaveConfig', {
          configType: this.configType,
          configName: this.configName,
          clientId: this.clientId || '',
          content: content
        }).subscribe({
          next: () => {
            this.isSaving = false;
            this.snackBar.open('Configuration saved successfully', 'Dismiss', { duration: 3000 });
          },
          error: (err) => {
            this.isSaving = false;
            this.error = `Failed to save configuration: ${err.message}`;
            this.snackBar.open(`Save failed: ${err.message}`, 'Dismiss', { duration: 5000 });
          }
        });
        */
    }

    setLanguageByConfigType(): void {
        switch (this.configType) {
            case 'ipxe':
                this.configLanguage = 'shell';
                break;
            case 'cloud-init':
                if (this.configName.includes('network')) {
                    this.configLanguage = 'yaml';
                } else if (this.configName === 'meta-data') {
                    this.configLanguage = 'json';
                } else {
                    this.configLanguage = 'yaml';
                }
                break;
            default:
                this.configLanguage = 'yaml';
        }
    }

    onConfigTypeChange(): void {
        this.loadAvailableConfigs();
    }

    onEditorInit(editor: any): void {
        // Additional editor configuration if needed
        editor.updateOptions({
            wordWrap: 'on',
            automaticLayout: true
        });
    }
}
