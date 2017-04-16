import {Component, Input, OnInit} from "@angular/core";
import {Modpack} from "app/models/modpack";

@Component({
  selector: 'app-file-upload',
  templateUrl: './file-upload.component.html',
  styleUrls: ['./file-upload.component.scss']
})
export class FileUploadComponent implements OnInit {

  @Input()
  protected modpack: Modpack;

  constructor() {
  }

  ngOnInit() {
  }

}
