import {TestBed, inject, async} from '@angular/core/testing';

import { ModpackService } from './modpack.service';

describe('ModpackService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ModpackService]
    });
  });

  it('should ...', inject([ModpackService], (service: ModpackService) => {
    expect(service).toBeTruthy();
  }));

  it('should add and remove modpacks', inject([ModpackService], (service: ModpackService) => {
    service.modpacks.take(1).subscribe(modpacks => {
      expect(modpacks.length).toBe(0);
    });
    const pack = service.addModpack("Test pack");
    service.modpacks.take(1).subscribe(modpacks => {
      if(expect(modpacks.length).toBe(1)) {
        expect(modpacks[0].name).toBe("Test pack");
      }
    });
    service.removeModpack(pack.id);
    service.modpacks.take(1).subscribe(modpacks => {
      expect(modpacks.length).toBe(0);
    });
  }));

  it('should set the selected modpack', inject([ModpackService], (service: ModpackService) => {
    const pack = service.addModpack("Test pack");
    service.setSelectedModpack(pack.id);
    service.selectedModpack.take(1).subscribe(modpack => {
      expect(modpack.id).toBe(pack.id);
    });
    const pack2 = service.addModpack("Test pack 2");
    service.setSelectedModpack(pack2.id);
    service.selectedModpack.take(1).subscribe(modpack => {
      expect(modpack.id).toBe(pack2.id);
    });
  }));
});
