import {ChangeDetectionStrategy, Component, Input, OnInit} from "@angular/core";
import {MdSnackBar} from "@angular/material";
import {Modpack} from "app/models/modpack";
import {BackendCommunicationService} from "app/services/backend-communication.service";


@Component({
  selector: 'app-solder-options',
  templateUrl: './solder-options.component.html',
  styleUrls: ['./solder-options.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class SolderOptionsComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  constructor(protected snackBar: MdSnackBar, protected backendCommunicationService: BackendCommunicationService) {
  }

  ngOnInit() {
    this.backendCommunicationService.getMessages<string>('solder-test')
      .subscribe(message => {
        this.snackBar.open(message, null, {duration: 5000});
      });
  }

  public testSolder() {
    this.backendCommunicationService.send('test-solder', this.modpack.solder);
  }
}
