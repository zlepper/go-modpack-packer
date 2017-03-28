import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {NoModpackSelectedComponent} from "./no-modpack-selected.component";

describe('NoModpackSelectedComponent', () => {
  let component: NoModpackSelectedComponent;
  let fixture: ComponentFixture<NoModpackSelectedComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [NoModpackSelectedComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NoModpackSelectedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
