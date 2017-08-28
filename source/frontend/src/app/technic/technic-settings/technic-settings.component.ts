import {Component, Input, OnInit} from "@angular/core";
import {MdDialog} from "@angular/material";
import {Modpack} from "app/models/modpack";
import {ModpackService} from "app/services/modpack.service";
import {TechnicBuildingComponent} from "app/technic/technic-building/technic-building.component";

@Component({
  selector: 'app-technic-settings',
  templateUrl: './technic-settings.component.html',
  styleUrls: ['./technic-settings.component.scss']
})
export class TechnicSettingsComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  constructor(protected dialog: MdDialog, private modpackService: ModpackService) {
  }

  ngOnInit() {
  }

  public startBuild() {
    // Save the modpack, so it's easy to restore the next time
    this.modpackService.saveModpacks();

    const dialog = this.dialog.open(TechnicBuildingComponent, {
      disableClose: true
    });
    dialog.afterClosed().subscribe(() => {
      console.log('closed');
    })
  }
}

