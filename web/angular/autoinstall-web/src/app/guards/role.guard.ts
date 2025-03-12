// filepath: /web/angular/src/app/guards/role.guard.ts
import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

@Injectable({
    providedIn: 'root'
})
export class RoleGuard implements CanActivate {
    constructor(private authService: AuthService, private router: Router) { }

    canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
        const requiredRole = route.data['requiredRole'];

        if (!requiredRole) {
            return true;
        }

        const hasRole = this.authService.hasRole(requiredRole);

        if (!hasRole) {
            this.router.navigate(['/access-denied']);
            return false;
        }

        return true;
    }
}
