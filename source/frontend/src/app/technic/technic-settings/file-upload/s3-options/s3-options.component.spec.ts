import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {S3OptionsComponent} from "./s3-options.component";

describe('S3OptionsComponent', () => {
  let component: S3OptionsComponent;
  let fixture: ComponentFixture<S3OptionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [S3OptionsComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(S3OptionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
