import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ModpackHeaderComponent } from './modpack-header.component';

describe('ModpackHeaderComponent', () => {
  let component: ModpackHeaderComponent;
  let fixture: ComponentFixture<ModpackHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ModpackHeaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModpackHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
