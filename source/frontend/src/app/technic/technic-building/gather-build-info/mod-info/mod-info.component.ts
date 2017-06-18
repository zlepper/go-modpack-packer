import {Component, Input, OnInit} from "@angular/core";
import {Mod} from "app/models/mod";
import {Modpack} from "app/models/modpack";

@Component({
  selector: 'app-mod-info',
  templateUrl: './mod-info.component.html',
  styleUrls: ['./mod-info.component.scss']
})
export class ModInfoComponent implements OnInit {
  public showDetails: boolean;

  @Input()
  public mod: Mod;

  @Input()
  public modpack: Modpack;

  constructor() { }

  ngOnInit() {
    this.showDetails = false;
  }

}
