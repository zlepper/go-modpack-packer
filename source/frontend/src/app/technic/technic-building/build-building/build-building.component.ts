import {Component, Input, OnInit} from '@angular/core';
import {Modpack, UploadWaiting} from "app/models/modpack";

@Component({
  selector: 'app-build-building',
  templateUrl: './build-building.component.html',
  styleUrls: ['./build-building.component.scss']
})
export class BuildBuildingComponent implements OnInit {

  @Input()
  protected modpack: Modpack;

  @Input()
  protected progressNumber: number;

  @Input()
  protected total: number;

  @Input()
  protected uploadNumber: number;

  @Input()
  protected uploading: string;

  @Input()
  protected uploadData: UploadWaiting;

  @Input()
  protected solderNumber: number;

  constructor() { }


  ngOnInit() {
  }

}
