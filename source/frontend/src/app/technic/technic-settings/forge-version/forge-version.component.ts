import {Component, Input, OnChanges, OnInit, SimpleChanges} from "@angular/core";
import {ForgeVersion} from "app/models/forgeversion";
import {Modpack} from "app/models/modpack";
import {ForgeVersionService} from "app/services/forge-version.service";
import {Observable} from "rxjs/Observable";

@Component({
  selector: 'app-forge-version',
  templateUrl: './forge-version.component.html',
  styleUrls: ['./forge-version.component.scss']
})
export class ForgeVersionComponent implements OnInit, OnChanges {

  @Input()
  public modpack: Modpack;
  public forgeVersions: Observable<ForgeVersion[]>;
  public ready: Observable<boolean>;

  constructor(protected forgeVersionService: ForgeVersionService) {
  }

  ngOnInit() {
    this.ready = this.forgeVersionService.ready;
  }

  ngOnChanges(changes: SimpleChanges): void {
    this.forgeVersions = this.forgeVersionService.getForgeVersionsForMCVersion(this.modpack.minecraftVersion);
  }
}
