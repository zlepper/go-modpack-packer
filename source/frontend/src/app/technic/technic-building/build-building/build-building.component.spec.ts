import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BuildBuildingComponent } from './build-building.component';

describe('BuildBuildingComponent', () => {
  let component: BuildBuildingComponent;
  let fixture: ComponentFixture<BuildBuildingComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BuildBuildingComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BuildBuildingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
