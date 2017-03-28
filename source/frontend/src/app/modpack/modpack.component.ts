import {ChangeDetectionStrategy, Component, OnInit} from "@angular/core";
import {ModpackService} from "app/services/modpack.service";
import {Observable} from "rxjs";
import {Modpack} from "app/models/modpack";

@Component({
  selector: 'app-modpack',
  templateUrl: './modpack.component.html',
  styleUrls: ['./modpack.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ModpackComponent implements OnInit {
  protected selectedModpack: Observable<Modpack>;

  constructor(protected modpackService: ModpackService) {
  }

  ngOnInit() {
    this.selectedModpack = this.modpackService.selectedModpack;
  }

}
