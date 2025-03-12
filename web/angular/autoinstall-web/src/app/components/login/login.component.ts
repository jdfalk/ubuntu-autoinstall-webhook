// filepath: /web/angular/autoinstall-web/src/app/components/login/login.component.ts
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { AuthService, AuthConfig } from '../../services/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Observable } from 'rxjs';

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
    loginForm: FormGroup;
    isLoading = false;
    authMethods: string[] = [];
    selectedAuthMethod = 'static';
    authConfig: AuthConfig | null = null;

    constructor(
        private fb: FormBuilder,
        private authService: AuthService,
        private router: Router,
        private route: ActivatedRoute,
        private snackBar: MatSnackBar
    ) {
        this.loginForm = this.fb.group({
            username: ['', [Validators.required]],
            password: ['', [Validators.required]]
        });
    }

    ngOnInit(): void {
        // Check if already authenticated
        if (this.authService.isAuthenticated()) {
            this.router.navigate(['/dashboard']);
            return;
        }

        // Check for OAuth callback
        this.route.queryParams.subscribe(params => {
            if (params['code']) {
                this.isLoading = true;
                this.authService.handleOAuthCallback(params['code']).subscribe({
                    next: () => {
                        this.router.navigate(['/dashboard']);
                    },
                    error: (err) => {
                        this.snackBar.open(err.message || 'OAuth authentication failed', 'Dismiss', {
                            duration: 5000
                        });
                        this.isLoading = false;
                    }
                });
            }
        });

        // Load available auth methods
        this.authService.getAuthConfig().subscribe(config => {
            if (config) {
                this.authConfig = config;
                this.authMethods = [];

                // Add enabled auth methods
                if (config.methods.static.enabled) {
                    this.authMethods.push('static');
                }

                if (config.methods.oauth.enabled) {
                    this.authMethods.push('oauth');
                }

                if (config.methods.database.enabled) {
                    this.authMethods.push('database');
                }

                if (config.methods.ldap.enabled) {
                    this.authMethods.push('ldap');
                }

                // Set default auth method
                if (this.authMethods.length > 0) {
                    this.selectedAuthMethod = this.authMethods[0];
                    this.onAuthMethodChange(this.selectedAuthMethod);
                }
            }
        });
    }

    onSubmit(): void {
        if (this.loginForm.invalid) {
            return;
        }

        this.isLoading = true;
        const username = this.loginForm.get('username')?.value;
        const password = this.loginForm.get('password')?.value;

        this.authService.login(username, password, this.selectedAuthMethod)
            .subscribe({
                next: () => {
                    this.router.navigate(['/dashboard']);
                },
                error: (err) => {
                    this.snackBar.open(err.message || 'Authentication failed', 'Dismiss', {
                        duration: 5000
                    });
                    this.isLoading = false;
                }
            });
    }

    initiateOAuthLogin(): void {
        const provider = this.authConfig?.methods.oauth.provider_type || 'generic';
        this.authService.initiateOAuthLogin(provider);
    }

    onAuthMethodChange(method: string): void {
        this.selectedAuthMethod = method;
        // Adjust form validation based on method
        if (method === 'oauth') {
            this.loginForm.get('password')?.clearValidators();
            this.loginForm.get('password')?.updateValueAndValidity();
        } else {
            this.loginForm.get('password')?.setValidators([Validators.required]);
            this.loginForm.get('password')?.updateValueAndValidity();
        }
    }

    getAuthMethodLabel(method: string): string {
        switch (method) {
            case 'static':
                return 'Static Credentials';
            case 'oauth':
                return 'OAuth / OIDC';
            case 'database':
                return 'Database User';
            case 'ldap':
                return 'LDAP / Active Directory';
            default:
                return method;
        }
    }
}
