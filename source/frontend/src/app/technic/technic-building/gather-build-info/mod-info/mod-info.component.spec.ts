import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ModInfoComponent } from './mod-info.component';

describe('ModInfoComponent', () => {
  let component: ModInfoComponent;
  let fixture: ComponentFixture<ModInfoComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ModInfoComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
