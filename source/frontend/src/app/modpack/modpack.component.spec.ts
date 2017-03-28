import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {ModpackComponent} from "./modpack.component";

describe('ModpackComponent', () => {
  let component: ModpackComponent;
  let fixture: ComponentFixture<ModpackComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ModpackComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModpackComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
