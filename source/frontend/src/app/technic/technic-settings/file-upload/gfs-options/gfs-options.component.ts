import {Component, Input, OnInit} from '@angular/core';
import {MdSnackBar} from "@angular/material";
import {Modpack} from "app/models/modpack";
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {BehaviorSubject} from "rxjs/BehaviorSubject";
import {Subject} from "rxjs/Subject";

export interface IGfsTestingResult {
  message: string;
}

@Component({
  selector: 'app-gfs-options',
  templateUrl: './gfs-options.component.html',
  styleUrls: ['./gfs-options.component.scss']
})
export class GfsOptionsComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  public testing: Subject<boolean>;

  constructor(private backendCommunicationService: BackendCommunicationService, private snackBar: MdSnackBar) {
  }

  ngOnInit() {
    this.testing = new BehaviorSubject<boolean>(false);
    this.backendCommunicationService.getMessages<IGfsTestingResult>("gfs-test")
      .subscribe(result => {
        this.snackBar.open(result.message, '', {
          duration: 5000
        });
        this.testing.next(false);
      });
  }

  public testGfs() {
    this.testing.next(true);
    this.backendCommunicationService.send('test-gfs', this.modpack.technic.upload.gfs);
  }

}
