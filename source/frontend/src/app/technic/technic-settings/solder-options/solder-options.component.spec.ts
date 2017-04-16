import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {SolderOptionsComponent} from "./solder-options.component";

describe('SolderOptionsComponent', () => {
  let component: SolderOptionsComponent;
  let fixture: ComponentFixture<SolderOptionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [SolderOptionsComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SolderOptionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
