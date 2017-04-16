import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {TechnicCheckPermissionsComponent} from "./technic-check-permissions.component";

describe('TechnicCheckPermissionsComponent', () => {
  let component: TechnicCheckPermissionsComponent;
  let fixture: ComponentFixture<TechnicCheckPermissionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TechnicCheckPermissionsComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TechnicCheckPermissionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
