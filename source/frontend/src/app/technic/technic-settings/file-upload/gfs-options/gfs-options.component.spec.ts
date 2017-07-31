import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {GfsOptionsComponent} from './gfs-options.component';

describe('GfsOptionsComponent', () => {
  let component: GfsOptionsComponent;
  let fixture: ComponentFixture<GfsOptionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [GfsOptionsComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GfsOptionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
