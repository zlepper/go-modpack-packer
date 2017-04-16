import {inject, TestBed} from "@angular/core/testing";

import {BackendCommunicationService} from "./backend-communication.service";

describe('BackendCommunicationService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [BackendCommunicationService]
    });
  });

  it('should ...', inject([BackendCommunicationService], (service: BackendCommunicationService) => {
    expect(service).toBeTruthy();
  }));
});
