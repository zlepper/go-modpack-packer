import {Component, Input, OnInit} from "@angular/core";
import {Mod} from "app/models/mod";
import {Modpack} from "app/models/modpack";
import {ModpackService} from "app/services/modpack.service";
import {Observable} from "rxjs/Observable";

@Component({
  selector: 'app-gather-build-info',
  templateUrl: './gather-build-info.component.html',
  styleUrls: ['./gather-build-info.component.scss']
})
export class GatherBuildInfoComponent implements OnInit {

  @Input()
  public mods: Observable<Mod[]>;

  @Input()
  public totalToScan: number;

  @Input()
  public showDone: Observable<boolean>;

  public modpack: Observable<Modpack>;

  public modsToShow: Observable<Mod[]>;

  constructor(protected modpackService: ModpackService) {
  }

  ngOnInit() {
    this.modsToShow = this.mods
      .withLatestFrom(this.showDone)
      .map(([mods, showDone]) => mods.filter(mod => showDone || !mod.isValid()));

    this.modpack = this.modpackService.selectedModpack;
  }


}
