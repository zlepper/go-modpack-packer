import {Component, Input, OnInit} from "@angular/core";
import {Mod} from "app/models/mod";
import {UserPermission} from "app/models/modpack";
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {ModpackService} from "app/services/modpack.service";


class PermissionSearch {
  public modId: string;
  public isPublic: boolean;

  constructor(id: string, isPublic: boolean) {
    this.modId = id;
    this.isPublic = isPublic;
  }
}


@Component({
  selector: 'app-technic-mod-permissions',
  templateUrl: './technic-mod-permissions.component.html',
  styleUrls: ['./technic-mod-permissions.component.scss']
})
export class TechnicModPermissionsComponent implements OnInit {

  @Input()
  public mod: Mod;

  constructor(protected backendCommunicationService: BackendCommunicationService, protected modpackService: ModpackService) {
  }

  ngOnInit() {
  }

  public checkDBForPermissions() {
    this.modpackService.selectedModpack
      .take(1)
      .switchMap(modpack => {
        this.backendCommunicationService.send('check-permissions-store', new PermissionSearch(this.mod.modid, modpack.technic.isPublicPack));
        return this.backendCommunicationService.getMessages<UserPermission>('get-permission-data');
      })
      .filter(permission => permission.modId === this.mod.modid)
      .take(1)
      .subscribe(permission => this.mod.userPermission = permission);
  }

}
