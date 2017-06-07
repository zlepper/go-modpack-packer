import {Component, Input, OnInit} from "@angular/core";
import {Modpack} from "app/models/modpack";
import {MdDialog} from '@angular/material';
import {TechnicBuildingComponent} from "app/technic/technic-building/technic-building.component";

@Component({
  selector: 'app-technic-settings',
  templateUrl: './technic-settings.component.html',
  styleUrls: ['./technic-settings.component.scss']
})
export class TechnicSettingsComponent implements OnInit {

  @Input()
  protected modpack: Modpack;

  constructor(protected dialog: MdDialog) {
  }

  ngOnInit() {
  }

  startBuild() {
    const dialog = this.dialog.open(TechnicBuildingComponent);
    dialog.afterClosed().subscribe(() => {
      console.log('closed');
    })
  }
}

