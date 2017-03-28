import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {TechnicComponent} from "./technic.component";

describe('TechnicComponent', () => {
  let component: TechnicComponent;
  let fixture: ComponentFixture<TechnicComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TechnicComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TechnicComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
