import {Component, Input, OnInit} from "@angular/core";
import {Modpack, UploadWaiting} from "app/models/modpack";

@Component({
  selector: 'app-build-building',
  templateUrl: './build-building.component.html',
  styleUrls: ['./build-building.component.scss']
})
export class BuildBuildingComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  @Input()
  public progressNumber: number;

  @Input()
  public total: number;

  @Input()
  public uploadNumber: number;

  @Input()
  public uploading: string;

  @Input()
  public uploadData: UploadWaiting;

  @Input()
  public solderNumber: number;

  constructor() { }


  ngOnInit() {
  }

}
