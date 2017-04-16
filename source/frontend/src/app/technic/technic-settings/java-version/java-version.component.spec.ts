import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {JavaVersionComponent} from "./java-version.component";

describe('JavaVersionComponent', () => {
  let component: JavaVersionComponent;
  let fixture: ComponentFixture<JavaVersionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [JavaVersionComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(JavaVersionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
