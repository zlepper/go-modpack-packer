import {Component, Input, OnInit} from "@angular/core";
import {MdDialog} from "@angular/material";
import {Modpack} from "app/models/modpack";
import {TechnicBuildingComponent} from "app/technic/technic-building/technic-building.component";

@Component({
  selector: 'app-technic-settings',
  templateUrl: './technic-settings.component.html',
  styleUrls: ['./technic-settings.component.scss']
})
export class TechnicSettingsComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  constructor(protected dialog: MdDialog) {
  }

  ngOnInit() {
  }

  public startBuild() {
    const dialog = this.dialog.open(TechnicBuildingComponent);
    dialog.afterClosed().subscribe(() => {
      console.log('closed');
    })
  }
}

