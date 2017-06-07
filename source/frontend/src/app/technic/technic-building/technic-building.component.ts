import {ChangeDetectionStrategy, Component, OnInit} from '@angular/core';
import {Modpack, UploadWaiting, UserPermission} from "app/models/modpack";
import {Observable} from "rxjs/Observable";
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {ModpackService} from "app/services/modpack.service";
import {Mod} from "app/models/mod";
import {Subject} from "rxjs/Subject";
import {BehaviorSubject} from "rxjs/BehaviorSubject";
import {MdDialogRef, MdCheckboxChange, MdSnackBar} from '@angular/material';
import 'app/operators/behaviorSubject';



@Component({
  selector: 'app-technic-building',
  templateUrl: './technic-building.component.html',
  styleUrls: ['./technic-building.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class TechnicBuildingComponent implements OnInit {

  protected modpack: Observable<Modpack>;
  protected mods: Observable<Mod[]>;
  protected readyToBuild: Observable<boolean>;
  protected totalToScan: Observable<number>;
  protected totalToPack: Observable<number>;
  protected packingTodos: Observable<string[]>;
  protected packingProgressNumber: Observable<number>;
  protected state: Subject<string>;
  protected showDone: Subject<boolean>;
  protected uploading: Observable<string>;
  protected updateTodos: Observable<string[]>;
  protected solderDoing: Observable<string>;
  protected uploadData: Observable<UploadWaiting>;
  protected solderProgressNumber: Observable<number>;
  protected uploadProgressNumber: Observable<number>;

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

  skipAll() {
    this.mods.take(1).subscribe(mods => mods.filter(mod => !mod.isValid()).forEach(mod => mod.skip = true));
  }

  changeShowDone(event: MdCheckboxChange) {
    this.showDone.next(event.checked)
  }

  build() {
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

  cancel() {
    this.dialogRef.close();
  }
}
