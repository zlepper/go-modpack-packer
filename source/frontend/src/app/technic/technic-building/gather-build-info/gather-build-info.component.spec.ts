import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GatherBuildInfoComponent } from './gather-build-info.component';

describe('GatherBuildInfoComponent', () => {
  let component: GatherBuildInfoComponent;
  let fixture: ComponentFixture<GatherBuildInfoComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ GatherBuildInfoComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GatherBuildInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
