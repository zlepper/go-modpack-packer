import {Component, Input, OnInit} from "@angular/core";
import {Modpack} from "app/models/modpack";

@Component({
  selector: 'app-java-version',
  templateUrl: './java-version.component.html',
  styleUrls: ['./java-version.component.scss']
})
export class JavaVersionComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  constructor() {
  }

  ngOnInit() {
  }

}
