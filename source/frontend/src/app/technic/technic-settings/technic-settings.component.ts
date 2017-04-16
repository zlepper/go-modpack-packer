import {Component, Input, OnInit} from "@angular/core";
import {Modpack} from "app/models/modpack";

@Component({
  selector: 'app-technic-settings',
  templateUrl: './technic-settings.component.html',
  styleUrls: ['./technic-settings.component.scss']
})
export class TechnicSettingsComponent implements OnInit {

  @Input()
  protected modpack: Modpack;

  constructor() {
  }

  ngOnInit() {
  }

}

