import {Component, Input, OnInit} from '@angular/core';
import {Mod} from "app/models/mod";
import {Observable} from "rxjs/Observable";
import {Modpack} from "app/models/modpack";
import {ModpackService} from "app/services/modpack.service";

@Component({
  selector: 'app-gather-build-info',
  templateUrl: './gather-build-info.component.html',
  styleUrls: ['./gather-build-info.component.scss']
})
export class GatherBuildInfoComponent implements OnInit {

  @Input()
  protected mods: Observable<Mod[]>;

  @Input()
  protected totalToScan: number;

  @Input()
  protected showDone: Observable<boolean>;

  protected modpack: Observable<Modpack>;

  protected modsToShow: Observable<Mod[]>;

  constructor(protected modpackService: ModpackService) {
  }

  ngOnInit() {
    this.modsToShow = this.mods
      .withLatestFrom(this.showDone)
      .map(([mods, showDone]) => mods.filter(mod => showDone || !mod.isValid()));

    this.modpack = this.modpackService.selectedModpack;
  }


}
