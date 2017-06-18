import {Component, OnInit} from "@angular/core";
import {Modpack} from "app/models/modpack";
import {ModpackService} from "app/services/modpack.service";
import {Observable} from "rxjs";

@Component({
  selector: 'app-modpack',
  templateUrl: './modpack.component.html',
  styleUrls: ['./modpack.component.scss']
})
export class ModpackComponent implements OnInit {
  public selectedModpack: Observable<Modpack>;

  constructor(protected modpackService: ModpackService) {
  }

  ngOnInit() {
    this.selectedModpack = this.modpackService.selectedModpack;
  }

}
