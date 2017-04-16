import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {FtpOptionsComponent} from "./ftp-options.component";

describe('FtpOptionsComponent', () => {
  let component: FtpOptionsComponent;
  let fixture: ComponentFixture<FtpOptionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [FtpOptionsComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FtpOptionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
