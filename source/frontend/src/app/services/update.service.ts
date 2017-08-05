import {Injectable} from '@angular/core';
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {ModpackService} from "app/services/modpack.service";
import {BehaviorSubject} from "rxjs/BehaviorSubject";
import {Observable} from "rxjs/Observable";
import {Subject} from "rxjs/Subject";

@Injectable()
export class UpdateService {

  constructor(private backendCommunicationService: BackendCommunicationService, private ModpackService: ModpackService) {
    this._updateAvailable = backendCommunicationService.getMessages('new-version-available')
      .map(() => true)
      .behaviorSubject(false);

    this._updating = new BehaviorSubject(false);

    this._updateMessage = this.backendCommunicationService.getMessages('update-progress');

    this.backendCommunicationService.getMessages('reload-frontend')
      .delay(5000)
      .subscribe(() => {
        location.reload(true);
      });
  }

  private _updateAvailable: Observable<boolean>;

  public get updateAvailable(): Observable<boolean> {
    return this._updateAvailable;
  }

  private _updating: Subject<boolean>;

  public get updating(): Observable<boolean> {
    return this._updating;
  }

  private _updateMessage: Observable<string>;

  public get updateMessage(): Observable<string> {
    return this._updateMessage;
  }

  public startUpdate() {
    this._updating.next(true);
    this.ModpackService.saveModpacks();
    this.backendCommunicationService.send('start-update', "")
  }
}
