import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {TechnicSettingsComponent} from "./technic-settings.component";

describe('TechnicSettingsComponent', () => {
  let component: TechnicSettingsComponent;
  let fixture: ComponentFixture<TechnicSettingsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TechnicSettingsComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TechnicSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
