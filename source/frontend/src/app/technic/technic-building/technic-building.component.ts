import {ChangeDetectionStrategy, Component, OnInit} from "@angular/core";
import {MdCheckboxChange, MdDialogRef, MdSnackBar} from "@angular/material";
import {Mod} from "app/models/mod";
import {Modpack, UploadWaiting} from "app/models/modpack";
import "app/operators/behaviorSubject";
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {ModpackService} from "app/services/modpack.service";
import {BehaviorSubject} from "rxjs/BehaviorSubject";
import {Observable} from "rxjs/Observable";
import {Subject} from "rxjs/Subject";


@Component({
  selector: 'app-technic-building',
  templateUrl: './technic-building.component.html',
  styleUrls: ['./technic-building.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class TechnicBuildingComponent implements OnInit {

  public modpack: Observable<Modpack>;
  public mods: Observable<Mod[]>;
  public readyToBuild: Observable<boolean>;
  public totalToScan: Observable<number>;
  public totalToPack: Observable<number>;
  public packingTodos: Observable<string[]>;
  public packingProgressNumber: Observable<number>;
  public state: Subject<string>;
  public showDone: Subject<boolean>;
  public uploading: Observable<string>;
  public updateTodos: Observable<string[]>;
  public solderDoing: Observable<string>;
  public uploadData: Observable<UploadWaiting>;
  public solderProgressNumber: Observable<number>;
  public uploadProgressNumber: Observable<number>;

  constructor(protected backendCommunicationService: BackendCommunicationService,
              protected modpackService: ModpackService,
              protected dialogRef: MdDialogRef<TechnicBuildingComponent>,
              protected snackBar: MdSnackBar) {
  }

  ngOnInit() {
    this.modpack = this.modpackService.selectedModpack;
    this.mods = new BehaviorSubject<Mod[]>([]);
    this.state = new BehaviorSubject('info');
    this.showDone = new BehaviorSubject(false);

    this.modpack.take(1).subscribe(modpack => {
      this.backendCommunicationService.send('gather-information', modpack);
    });

    this.mods = this.backendCommunicationService.getMessages<Mod>('mod-data-ready')
      .map(mod => Mod.fromJson(mod))
      .bufferTime(10)
      .scan((currentMods: Mod[], newMods: Mod[]) => [...currentMods, ...newMods], [])
      .behaviorSubject([]);

    this.readyToBuild = this.backendCommunicationService.getMessages('all-mod-files-scanned').behaviorSubject(false);

    this.totalToScan = this.backendCommunicationService.getMessages<number>('total-mod-files').behaviorSubject(-1);

    this.totalToPack = this.backendCommunicationService.getMessages<number>('total-to-pack').behaviorSubject(-1);

    const donePackingPart = this.backendCommunicationService.getMessages<string>('done-packing-part')
      .bufferTime(10);

    this.packingTodos = this.backendCommunicationService.getMessages<string>('packing-part')
      .bufferTime(10)
      .withLatestFrom(donePackingPart)
      .scan((currentTodos, [newTodos, todosToRemove]) =>
          [...currentTodos, ...newTodos]
            .filter((todo: string) => todosToRemove.includes(todo))
        , []);

    this.packingProgressNumber = donePackingPart.scan((currentProgress, newTodos) => currentProgress + newTodos.length, 0);

    this.uploading = this.backendCommunicationService.getMessages<string>('starting-upload')
      .withLatestFrom(this.backendCommunicationService.getMessages<string>('finished-all-uploading').mapTo(true).behaviorSubject(false))
      .map(([uploading, done]) => done ? '' : uploading);

    this.uploading.subscribe(uploading => console.log(uploading));

    this.uploadProgressNumber = this.backendCommunicationService.getMessages<string>('finished-uploading').scan(currentProgress => currentProgress + 1, 0);

    const doneUpdatingSolder = this.backendCommunicationService.getMessages<string>('done-updating-solder').bufferTime(10);

    this.updateTodos = this.backendCommunicationService.getMessages<string>('updating-solder')
      .bufferTime(10)
      .withLatestFrom(doneUpdatingSolder)
      .scan((currentTodos, [newTodos, todosToRemove]) =>
          [...currentTodos, ...newTodos]
            .filter((todo: string) => todosToRemove.includes(todo))
        , []);

    this.solderDoing = this.backendCommunicationService.getMessages<string>('solder-currently-doing').behaviorSubject('');

    this.solderDoing.subscribe(doing => console.log(doing));

    this.solderProgressNumber = doneUpdatingSolder.scan((currentProgress, doneTodos) => currentProgress + doneTodos.length, 0);

    this.uploadData = this.backendCommunicationService.getMessages<UploadWaiting>('waiting-for-file-upload');

  }

  public skipAll() {
    this.mods.take(1).subscribe(mods => mods.filter(mod => !mod.isValid()).forEach(mod => mod.skip = true));
  }

  public changeShowDone(event: MdCheckboxChange) {
    this.showDone.next(event.checked)
  }

  public build() {
    this.mods.take(1)
      .withLatestFrom(this.modpack.take(1))
      .subscribe(([mods, modpack]) => {
        const shouldBuild = mods.reduce((valid, mod) => valid && (mod.skip || mod.isValid()), true);
        if (shouldBuild) {
          mods = mods.filter(m => !m.skip);
          this.backendCommunicationService.send('build', {modpack, mods});
          this.state.next('building');
        } else {
          this.snackBar.open('Some mods are missing info. Please fill it in before continuing.', null, {duration: 5000});
        }
      });
  }

  public cancel() {
    this.dialogRef.close();
  }
}
