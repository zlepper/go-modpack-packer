import {Component, Input, OnInit} from "@angular/core";
import {Modpack} from "app/models/modpack";

@Component({
  selector: 'app-pack-type',
  templateUrl: './pack-type.component.html',
  styleUrls: ['./pack-type.component.scss']
})
export class PackTypeComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  constructor() {
  }

  ngOnInit() {
  }

}
