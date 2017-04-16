import {Component, OnInit} from "@angular/core";
import {Modpack} from "app/models/modpack";
import {ModpackService} from "app/services/modpack.service";
import {Observable} from "rxjs/Observable";

@Component({
  selector: 'app-technic',
  templateUrl: './technic.component.html',
  styleUrls: ['./technic.component.scss']
})
export class TechnicComponent implements OnInit {

  protected selectedModpack: Observable<Modpack>;

  constructor(protected modpackService: ModpackService) {
  }

  ngOnInit() {
    this.selectedModpack = this.modpackService.selectedModpack;
  }

}
