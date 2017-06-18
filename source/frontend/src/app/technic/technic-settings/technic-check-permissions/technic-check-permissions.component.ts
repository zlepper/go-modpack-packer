import {Component, Input, OnInit} from "@angular/core";
import {Modpack} from "app/models/modpack";

@Component({
  selector: 'app-technic-check-permissions',
  templateUrl: './technic-check-permissions.component.html',
  styleUrls: ['./technic-check-permissions.component.scss']
})
export class TechnicCheckPermissionsComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  constructor() {
  }

  ngOnInit() {
  }

}
