import {ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output} from "@angular/core";
import {Modpack, UploadWaiting} from "app/models/modpack";
import {BehaviorSubject} from "rxjs/BehaviorSubject";
import {Subject} from "rxjs/Subject";

@Component({
  selector: 'app-build-building',
  templateUrl: './build-building.component.html',
  styleUrls: ['./build-building.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class BuildBuildingComponent implements OnInit {

  @Input()
  public modpack: Modpack;

  @Input()
  public progressNumber: number;

  @Input()
  public total: number;

  @Input()
  public uploadNumber: number;

  @Input()
  public uploading: string;

  @Input()
  public uploadData: UploadWaiting;

  @Input()
  public solderNumber: number;

  @Output()
  public continueBuild = new EventEmitter<void>();

  public continueClicked: Subject<boolean> = new BehaviorSubject(false);

  constructor() { }


  ngOnInit() {
  }

  continueBuildingWithManualStorage() {
    this.continueBuild.emit();
    this.continueClicked.next(true);
  }

}
