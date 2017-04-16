import {inject, TestBed} from "@angular/core/testing";

import {HasModpackSelectedGuard} from "./has-modpack-selected.guard";

describe('HasModpackSelectedGuard', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [HasModpackSelectedGuard]
    });
  });

  it('should ...', inject([HasModpackSelectedGuard], (guard: HasModpackSelectedGuard) => {
    expect(guard).toBeTruthy();
  }));
});
