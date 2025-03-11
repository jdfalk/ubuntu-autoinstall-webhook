import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface StatusRequest {
    hostname: string;
    ip_address: string;
    progress: number;
    message: string;
}

export interface StatusResponse {
    acknowledged: boolean;
}

@Injectable({
    providedIn: 'root'
})
export class InstallService {
    private baseUrl = 'http://localhost:8080/v1/install';

    constructor(private http: HttpClient) { }

    reportStatus(status: StatusRequest): Observable<StatusResponse> {
        return this.http.post<StatusResponse>(`${this.baseUrl}/status`, status);
    }
}
