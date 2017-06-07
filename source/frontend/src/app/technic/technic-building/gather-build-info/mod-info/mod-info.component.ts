import {Component, Input, OnInit} from '@angular/core';
import {Mod} from "app/models/mod";
import {Modpack} from "app/models/modpack";

@Component({
  selector: 'app-mod-info',
  templateUrl: './mod-info.component.html',
  styleUrls: ['./mod-info.component.scss']
})
export class ModInfoComponent implements OnInit {
  protected showDetails: boolean;

  @Input()
  protected mod: Mod;

  @Input()
  protected modpack: Modpack;


  constructor() { }

  ngOnInit() {
    this.showDetails = false;
  }

}
