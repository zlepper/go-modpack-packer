import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {PackTypeComponent} from "./pack-type.component";

describe('PackTypeComponent', () => {
  let component: PackTypeComponent;
  let fixture: ComponentFixture<PackTypeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [PackTypeComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PackTypeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
