import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ModpackEditorComponent } from './modpack-editor.component';

describe('ModpackEditorComponent', () => {
  let component: ModpackEditorComponent;
  let fixture: ComponentFixture<ModpackEditorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ModpackEditorComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModpackEditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
