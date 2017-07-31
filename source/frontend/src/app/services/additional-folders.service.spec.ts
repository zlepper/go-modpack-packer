import {inject, TestBed} from '@angular/core/testing';

import {AdditionalFoldersService} from './additional-folders.service';

describe('AdditionalFoldersService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [AdditionalFoldersService]
    });
  });

  it('should be created', inject([AdditionalFoldersService], (service: AdditionalFoldersService) => {
    expect(service).toBeTruthy();
  }));
});
