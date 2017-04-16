import {Component, Input, OnInit} from "@angular/core";
import {Modpack} from "app/models/modpack";
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {BehaviorSubject} from "rxjs/BehaviorSubject";
import {Observable} from "rxjs/Observable";
import {Subject} from "rxjs/Subject";

@Component({
  selector: 'app-s3-options',
  templateUrl: './s3-options.component.html',
  styleUrls: ['./s3-options.component.scss']
})
export class S3OptionsComponent implements OnInit {

  @Input()
  protected modpack: Modpack;
  protected buckets: Subject<string[]>;
  protected hasBuckets: Observable<boolean>;

  constructor(protected backendCommunicationService: BackendCommunicationService) {
  }

  ngOnInit() {
    this.buckets = new BehaviorSubject([]);
    this.backendCommunicationService.getMessages<string[]>('found-aws-buckets')
      .subscribe(buckets => this.buckets.next(buckets));
    this.hasBuckets = this.buckets.map(buckets => buckets.length > 0);
  }

  getBuckets() {
    this.backendCommunicationService.send('get-aws-buckets', this.modpack);
  }

}
