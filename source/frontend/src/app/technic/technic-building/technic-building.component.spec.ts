import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TechnicBuildingComponent } from './technic-building.component';

describe('TechnicBuildingComponent', () => {
  let component: TechnicBuildingComponent;
  let fixture: ComponentFixture<TechnicBuildingComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TechnicBuildingComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TechnicBuildingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
