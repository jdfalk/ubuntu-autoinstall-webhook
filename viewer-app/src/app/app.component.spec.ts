import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { AppComponent } from './app.component';

describe('AppComponent', () => {
  let component: AppComponent;
  let fixture: ComponentFixture<AppComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      providers: [provideHttpClientTesting(), provideRouter([])],
      declarations: [AppComponent]
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(AppComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should fetch systems on init', () => {
    spyOn(component['http'], 'get').and.callThrough();
    component.ngOnInit();
    expect(component['http'].get).toHaveBeenCalledWith('/api/viewer');
  });

  it('should navigate to logs when an IP is clicked', () => {
    spyOn(component['router'], 'navigate');
    component.navigateToLogs('192.168.1.1');
    expect(component['router'].navigate).toHaveBeenCalledWith(['/viewer', '192.168.1.1']);
  });
});
