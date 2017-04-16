import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {ForgeVersionComponent} from "./forge-version.component";

describe('ForgeVersionComponent', () => {
  let component: ForgeVersionComponent;
  let fixture: ComponentFixture<ForgeVersionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ForgeVersionComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ForgeVersionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
