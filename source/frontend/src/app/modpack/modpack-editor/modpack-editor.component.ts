import {Component, Input, OnInit} from '@angular/core';
import {Modpack} from "app/models/modpack";
import {ForgeVersionService} from "app/services/forge-version.service";

@Component({
  selector: 'app-modpack-editor',
  templateUrl: './modpack-editor.component.html',
  styleUrls: ['./modpack-editor.component.scss']
})
export class ModpackEditorComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  constructor(protected forgeVersionService: ForgeVersionService) { }

  ngOnInit() {
  }

}
