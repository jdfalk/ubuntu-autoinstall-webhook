import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
    selector: 'app-access-denied',
    templateUrl: './access-denied.component.html',
    styleUrls: ['./access-denied.component.scss']
})
export class AccessDeniedComponent {
    currentUser: any = null;

    constructor(
        private router: Router,
        private authService: AuthService
    ) {
        this.authService.getCurrentUser().subscribe(user => {
            this.currentUser = user;
        });
    }

    goToDashboard(): void {
        this.router.navigate(['/dashboard']);
    }
}
