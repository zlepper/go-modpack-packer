import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {FtbComponent} from "./ftb.component";

describe('FtbComponent', () => {
  let component: FtbComponent;
  let fixture: ComponentFixture<FtbComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [FtbComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FtbComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
