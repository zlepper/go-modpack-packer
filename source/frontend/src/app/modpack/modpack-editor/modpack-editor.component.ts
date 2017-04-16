import {Component, Input, OnInit} from "@angular/core";
import {Modpack} from "app/models/modpack";
import {ElectronService} from "app/services/electron.service";
import {ForgeVersionService} from "app/services/forge-version.service";

@Component({
  selector: 'app-modpack-editor',
  templateUrl: './modpack-editor.component.html',
  styleUrls: ['./modpack-editor.component.scss']
})
export class ModpackEditorComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  constructor(protected forgeVersionService: ForgeVersionService, protected electron: ElectronService) {
    this.electron.on('selected-input-directory', (path: string) => {
      this.modpack.inputDirectory = path;
    });
    this.electron.on('selected-output-directory', (path: string) => {
      this.modpack.outputDirectory = path;
    });
  }

  ngOnInit() {
  }

  selectInputDirectory() {
    this.electron.send('open-input-directory-dialog', null);
  }

  selectOutputDirectory(): void {
    this.electron.send('open-output-directory-dialog', null);
  }

}
