import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ConfigEditorComponent } from './config-editor.component';

describe('ConfigEditorComponent', () => {
  let component: ConfigEditorComponent;
  let fixture: ComponentFixture<ConfigEditorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      // If you’re using a standalone component approach, or if you’re
      // using a module-based approach, ensure all needed imports are here.
      imports: [
        // The component itself
        ConfigEditorComponent
        // Possibly MatToolbarModule, MatButtonModule, FormsModule, etc. if testing those.
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(ConfigEditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
