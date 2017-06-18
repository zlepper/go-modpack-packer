import {ChangeDetectionStrategy, Component, OnInit} from "@angular/core";
import {ModpackService} from "app/services/modpack.service";

@Component({
  selector: 'app-body',
  templateUrl: './body.component.html',
  styleUrls: ['./body.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class BodyComponent implements OnInit {

  public selectedModpackId: number;

  constructor(public modpackService: ModpackService) {

  }

  ngOnInit() {
    this.modpackService.selectedModpack
      .filter(modpack => modpack != null)
      .subscribe(modpack => {
        if (modpack.id !== this.selectedModpackId) {
          this.selectedModpackId = modpack.id;
        }
      });
  }

  public selectedModpackChanged(id: number) {
    if (id === -1) {
      const pack = this.modpackService.addModpack("Unnamed modpack");
      id = this.selectedModpackId = pack.id;
    }

    this.modpackService.setSelectedModpack(id);
  }

}
