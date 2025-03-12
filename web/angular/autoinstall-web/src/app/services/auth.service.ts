// filepath: /web/angular/src/app/services/auth.service.ts
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, of, throwError } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { tap, catchError, switchMap } from 'rxjs/operators';

export interface Role {
    name: string;
    childRoles?: string[];
}

export interface User {
    username: string;
    roles: string[];
}

export interface AuthResponse {
    success: boolean;
    user?: User;
    token?: string;
    message?: string;
}

export interface AuthConfig {
    methods: {
        static: {
            enabled: boolean;
        };
        oauth: {
            enabled: boolean;
            provider_type: string;
        };
        database: {
            enabled: boolean;
        };
        ldap: {
            enabled: boolean;
        };
    };
}

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    private currentUserRoles = new BehaviorSubject<string[]>([]);
    private currentUser = new BehaviorSubject<User | null>(null);
    private roleDefinitions: Map<string, Role> = new Map();
    private authToken = new BehaviorSubject<string | null>(null);
    private authConfig = new BehaviorSubject<AuthConfig | null>(null);

    // For testing purposes - should be removed in production
    private testUsers = [
        { username: 'admin', password: 'admin123', roles: ['admin'] },
        { username: 'user', password: 'user123', roles: ['user'] }
    ];

    constructor(private http: HttpClient) {
        // Initialize role definitions
        this.roleDefinitions.set('admin', { name: 'admin', childRoles: ['logging', 'ide', 'configuration'] });
        this.roleDefinitions.set('logging', { name: 'logging' });
        this.roleDefinitions.set('ide', { name: 'ide' });
        this.roleDefinitions.set('configuration', { name: 'configuration' });
        this.roleDefinitions.set('user', { name: 'user', childRoles: ['logging'] });

        // Check for existing session
        this.loadStoredSession();

        // Load auth configuration from server
        this.loadAuthConfig();
    }

    private loadStoredSession(): void {
        const storedUser = localStorage.getItem('currentUser');
        const storedToken = localStorage.getItem('authToken');

        if (storedUser) {
            try {
                const user = JSON.parse(storedUser) as User;
                this.currentUser.next(user);
                this.currentUserRoles.next(user.roles);

                if (storedToken) {
                    this.authToken.next(storedToken);
                }
            } catch (e) {
                this.clearSession();
            }
        }
    }

    private loadAuthConfig(): void {
        this.http.get<AuthConfig>(`${environment.apiEndpoint}/auth/config`)
            .pipe(
                catchError(() => {
                    // Use default config if we can't load from server
                    const defaultConfig: AuthConfig = {
                        methods: {
                            static: { enabled: true },
                            oauth: { enabled: false, provider_type: 'generic' },
                            database: { enabled: false },
                            ldap: { enabled: false }
                        }
                    };
                    return of(defaultConfig);
                })
            )
            .subscribe(config => {
                this.authConfig.next(config);
            });
    }

    getAuthConfig(): Observable<AuthConfig | null> {
        return this.authConfig.asObservable();
    }

    login(username: string, password: string, authMethod: string = 'static'): Observable<AuthResponse> {
        // For development/testing with static auth enabled, use static users
        if (authMethod === 'static' && this.authConfig.value?.methods.static.enabled) {
            const user = this.testUsers.find(
                u => u.username === username && u.password === password
            );

            if (user) {
                const authResponse: AuthResponse = {
                    success: true,
                    user: {
                        username: user.username,
                        roles: user.roles
                    },
                    token: 'mock-jwt-token'
                };

                // Store user info
                this.setSession(authResponse);

                return of(authResponse);
            }

            return throwError(() => new Error('Invalid username or password'));
        }

        // For other auth methods, call the backend
        return this.http.post<AuthResponse>(`${environment.apiEndpoint}/auth/login`, {
            username,
            password,
            method: authMethod
        }).pipe(
            tap(response => {
                if (response.success && response.user) {
                    this.setSession(response);
                }
            }),
            catchError(error => {
                return throwError(() => new Error(error.error?.message || 'Authentication failed'));
            })
        );
    }

    initiateOAuthLogin(provider: string = 'generic'): void {
        // Get the current URL for the redirect_uri
        const redirectUri = `${window.location.origin}/auth/callback`;

        // Open the OAuth login page
        this.http.get<{ auth_url: string }>(`${environment.apiEndpoint}/auth/oauth/init`, {
            params: { provider, redirect_uri: redirectUri }
        }).subscribe(
            response => {
                // Redirect to the auth URL
                window.location.href = response.auth_url;
            },
            error => {
                console.error('OAuth initialization failed', error);
            }
        );
    }

    handleOAuthCallback(code: string): Observable<AuthResponse> {
        // Exchange the code for a token
        return this.http.post<AuthResponse>(`${environment.apiEndpoint}/auth/oauth/callback`, { code })
            .pipe(
                tap(response => {
                    if (response.success && response.user) {
                        this.setSession(response);
                    }
                }),
                catchError(error => {
                    return throwError(() => new Error(error.error?.message || 'OAuth authentication failed'));
                })
            );
    }

    private setSession(authResponse: AuthResponse): void {
        if (authResponse.user) {
            this.currentUser.next(authResponse.user);
            this.currentUserRoles.next(authResponse.user.roles);
            localStorage.setItem('currentUser', JSON.stringify(authResponse.user));
        }

        if (authResponse.token) {
            this.authToken.next(authResponse.token);
            localStorage.setItem('authToken', authResponse.token);
        }
    }

    private clearSession(): void {
        this.currentUser.next(null);
        this.currentUserRoles.next([]);
        this.authToken.next(null);
        localStorage.removeItem('currentUser');
        localStorage.removeItem('authToken');
    }

    logout(): void {
        // Call logout endpoint if we have a token
        const token = this.authToken.value;
        if (token) {
            this.http.post(`${environment.apiEndpoint}/auth/logout`, {})
                .pipe(
                    catchError(err => {
                        console.error('Logout error:', err);
                        return of(null);
                    })
                )
                .subscribe(() => {
                    this.clearSession();
                });
        } else {
            this.clearSession();
        }
    }

    getToken(): string | null {
        return this.authToken.value;
    }

    getUserRoles(): Observable<string[]> {
        return this.currentUserRoles.asObservable();
    }

    getCurrentUser(): Observable<User | null> {
        return this.currentUser.asObservable();
    }

    hasRole(roleName: string): boolean {
        const userRoles = this.currentUserRoles.value;

        // Direct role check
        if (userRoles.includes(roleName)) {
            return true;
        }

        // Check inherited roles
        return userRoles.some(role => {
            const roleDefinition = this.roleDefinitions.get(role);
            return roleDefinition?.childRoles?.includes(roleName) || false;
        });
    }

    isAuthenticated(): boolean {
        return this.currentUserRoles.value.length > 0;
    }
}
