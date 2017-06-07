import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TechnicModPermissionsComponent } from './technic-mod-permissions.component';

describe('TechnicModPermissionsComponent', () => {
  let component: TechnicModPermissionsComponent;
  let fixture: ComponentFixture<TechnicModPermissionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TechnicModPermissionsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TechnicModPermissionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
