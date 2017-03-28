import { TestBed, inject } from '@angular/core/testing';

import { ForgeVersionService } from './forge-version.service';

describe('ForgeVersionService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ForgeVersionService]
    });
  });

  it('should ...', inject([ForgeVersionService], (service: ForgeVersionService) => {
    expect(service).toBeTruthy();
  }));
});
