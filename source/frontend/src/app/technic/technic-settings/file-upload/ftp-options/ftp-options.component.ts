import {Component, Input, OnInit} from "@angular/core";
import {MdSnackBar} from "@angular/material";
import {Modpack} from "app/models/modpack";
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {BehaviorSubject} from "rxjs/BehaviorSubject";
import {Subject} from "rxjs/Subject";

export interface IFtpTestResult {
  success: boolean;
  message: string;
}

@Component({
  selector: 'app-ftp-options',
  templateUrl: './ftp-options.component.html',
  styleUrls: ['./ftp-options.component.scss']
})
export class FtpOptionsComponent implements OnInit {

  @Input()
  protected modpack: Modpack;
  protected testing: Subject<boolean>;

  constructor(protected backendCommunicationService: BackendCommunicationService, protected snackBar: MdSnackBar) {
  }

  ngOnInit() {
    this.testing = new BehaviorSubject<boolean>(false);
    this.backendCommunicationService.getMessages<IFtpTestResult>("ftp-test")
      .subscribe(result => {
        this.snackBar.open(result.message, null, {
          duration: 5000
        });
        this.testing.next(false);
      });
  }

  testFtp() {
    console.log('Testing ftp');
    this.testing.next(true);
    this.backendCommunicationService.send('test-ftp', this.modpack.technic.upload.ftp);
  }
}
