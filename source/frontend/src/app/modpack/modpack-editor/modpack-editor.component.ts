import {Component, Input, OnInit} from "@angular/core";
import {FormControl} from "@angular/forms";
import {Modpack} from "app/models/modpack";
import 'app/operators/log';
import {FolderService} from "app/services/folders.service";
import {ForgeVersionService} from "app/services/forge-version.service";
import {Observable} from "rxjs/Observable";

@Component({
  selector: 'app-modpack-editor',
  templateUrl: './modpack-editor.component.html',
  styleUrls: ['./modpack-editor.component.scss']
})
export class ModpackEditorComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  public inputControl = new FormControl();
  public filteredInputFolders: Observable<string[]>;

  public outputControl = new FormControl();
  public filteredOutputFolders: Observable<string[]>;

  constructor(protected forgeVersionService: ForgeVersionService, private folderService: FolderService) {
  }

  ngOnInit() {
    this.filteredInputFolders = this.setupFolderWatch(this.inputControl);
    this.filteredOutputFolders = this.setupFolderWatch(this.outputControl);
  }

  private setupFolderWatch(input: FormControl): Observable<Array<string>> {
    return input.valueChanges
      .startWith('')
      .throttleTime(350)
      .map((folder: string) => folder.replace('\\', '/'))
      .switchMap(folder => this.folderService.search(folder));
  }

}
