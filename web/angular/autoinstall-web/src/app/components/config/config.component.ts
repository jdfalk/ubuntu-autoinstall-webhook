import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormArray, FormControl, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { GrpcService } from '../../services/grpc.service';
import { TelemetryService } from '../../services/telemetry.service';

@Component({
    selector: 'app-config',
    templateUrl: './config.component.html',
    styleUrls: ['./config.component.scss']
})
export class ConfigComponent implements OnInit {
    configForm: FormGroup;
    sections: string[] = ['server', 'database', 'auth', 'paths'];
    activeSection: string = 'server';
    isLoading = true;
    isSaving = false;
    error = '';

    constructor(
        private fb: FormBuilder,
        private grpcService: GrpcService,
        private snackBar: MatSnackBar,
        private telemetryService: TelemetryService
    ) {
        this.configForm = this.createConfigForm();
    }

    ngOnInit(): void {
        this.loadConfig();
    }

    createConfigForm(): FormGroup {
        return this.fb.group({
            // Server settings
            serverHost: ['0.0.0.0', Validators.required],
            serverPort: [8080, [Validators.required, Validators.min(1), Validators.max(65535)]],
            enableHttps: [false],
            httpsCertPath: [''],
            httpsKeyPath: [''],

            // Database settings
            dbType: ['sqlite3'],
            sqlitePath: ['/var/lib/ubuntu-autoinstall-webhook/database.sqlite'],
            cockroachHost: ['localhost'],
            cockroachPort: [26257],
            cockroachDb: ['ubuntu_autoinstall'],
            cockroachUser: ['root'],
            cockroachSslMode: ['disable'],

            // Auth settings
            staticAuthEnabled: [true],
            staticUsers: this.fb.array([]),
            oauthEnabled: [false],
            oauthProviderType: ['generic'],
            oauthClientId: [''],
            oauthClientSecret: [''],
            oauthAuthUrl: [''],
            oauthTokenUrl: [''],
            oauthRoleAttribute: ['groups'],
            ldapEnabled: [false],
            ldapServerUrl: [''],
            ldapBindDn: [''],
            ldapBindPassword: [''],
            ldapSearchBase: [''],
            ldapUserFilter: ['(uid=%s)'],
            ldapGroupFilter: ['(memberUid=%s)'],

            // Path settings
            ipxeBootPath: ['/var/www/html/ipxe/boot/'],
            cloudInitPath: ['/var/www/html/cloud-init/']
        });
    }

    loadConfig(): void {
        this.isLoading = true;
        this.error = '';

        // In a real implementation, this would call your gRPC service
        // For now, using setTimeout to simulate loading
        setTimeout(() => {
            // Simulate loading default config
            this.configForm.patchValue({
                serverHost: '0.0.0.0',
                serverPort: 8080,
                dbType: 'sqlite3',
                sqlitePath: '/var/lib/ubuntu-autoinstall-webhook/database.sqlite',
                ipxeBootPath: '/var/www/html/ipxe/boot/',
                cloudInitPath: '/var/www/html/cloud-init/',
                staticAuthEnabled: true
            });

            // Add mock users
            const staticUsers = this.configForm.get('staticUsers') as FormArray;
            staticUsers.clear();

            staticUsers.push(this.createUserFormGroup('admin', 'admin123', ['admin']));
            staticUsers.push(this.createUserFormGroup('user', 'user123', ['user']));

            this.isLoading = false;
        }, 1000);

        // Real implementation would be:
        /*
        this.grpcService.callService('ConfigService', 'GetConfig', {}).subscribe({
          next: (response) => {
            // Parse and set form values from response
            this.configForm.patchValue({
              serverHost: response.server.host,
              serverPort: response.server.port,
              // ... other fields
            });

            // Handle arrays and nested objects
            const staticUsers = this.configForm.get('staticUsers') as FormArray;
            staticUsers.clear();

            if (response.auth?.methods?.static?.users) {
              response.auth.methods.static.users.forEach((user: any) => {
                staticUsers.push(this.createUserFormGroup(
                  user.username,
                  user.password,
                  user.roles
                ));
              });
            }

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
        if (this.configForm.invalid) {
            this.snackBar.open('Please correct the errors before saving', 'Dismiss', { duration: 5000 });
            return;
        }

        this.isSaving = true;

        // Convert form values to the expected format for your API
        const formValues = this.configForm.value;
        const config = {
            server: {
                host: formValues.serverHost,
                port: formValues.serverPort
            },
            https: {
                enabled: formValues.enableHttps,
                cert_path: formValues.httpsCertPath,
                key_path: formValues.httpsKeyPath
            },
            database: {
                type: formValues.dbType,
                sqlite3: {
                    path: formValues.sqlitePath
                },
                cockroachdb: {
                    host: formValues.cockroachHost,
                    port: formValues.cockroachPort,
                    database: formValues.cockroachDb,
                    user: formValues.cockroachUser,
                    ssl_mode: formValues.cockroachSslMode
                }
            },
            auth: {
                methods: {
                    static: {
                        enabled: formValues.staticAuthEnabled,
                        users: formValues.staticUsers
                    },
                    oauth: {
                        enabled: formValues.oauthEnabled,
                        provider_type: formValues.oauthProviderType,
                        client_id: formValues.oauthClientId,
                        client_secret: formValues.oauthClientSecret,
                        auth_url: formValues.oauthAuthUrl,
                        token_url: formValues.oauthTokenUrl,
                        role_attribute: formValues.oauthRoleAttribute
                    },
                    ldap: {
                        enabled: formValues.ldapEnabled,
                        server_url: formValues.ldapServerUrl,
                        bind_dn: formValues.ldapBindDn,
                        bind_password: formValues.ldapBindPassword,
                        search_base: formValues.ldapSearchBase,
                        user_filter: formValues.ldapUserFilter,
                        group_filter: formValues.ldapGroupFilter
                    }
                }
            },
            paths: {
                ipxe_boot: formValues.ipxeBootPath,
                cloud_init: formValues.cloudInitPath
            }
        };

        // For now, just simulate an API call
        setTimeout(() => {
            this.isSaving = false;
            this.snackBar.open('Configuration saved successfully', 'Dismiss', { duration: 3000 });
        }, 1500);

        // Real implementation would be:
        /*
        this.grpcService.callService('ConfigService', 'SaveConfig', config).subscribe({
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

    resetForm(): void {
        this.loadConfig();
    }

    setActiveSection(section: string): void {
        this.activeSection = section;
    }

    getControl(name: string): FormControl {
        return this.configForm.get(name) as FormControl;
    }

    get staticUsersArray(): FormArray {
        return this.configForm.get('staticUsers') as FormArray;
    }

    createUserFormGroup(username: string = '', password: string = '', roles: string[] = []): FormGroup {
        return this.fb.group({
            username: [username, Validators.required],
            password: [password, Validators.required],
            roles: [roles, Validators.required]
        });
    }

    addStaticUser(): void {
        const staticUsers = this.configForm.get('staticUsers') as FormArray;
        staticUsers.push(this.createUserFormGroup());
    }

    removeStaticUser(index: number): void {
        const staticUsers = this.configForm.get('staticUsers') as FormArray;
        staticUsers.removeAt(index);
    }
}
