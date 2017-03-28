import {Injectable} from "@angular/core";
import {BehaviorSubject, Observable, Subject} from "rxjs";

declare var electron;
const remote = electron.remote;

@Injectable()
export class ElectronService {
  private _isMaximized: Subject<boolean>;

  constructor() {
    this._isMaximized = new BehaviorSubject<boolean>(remote.getCurrentWindow().isMaximized());
  }

  public get isMaximized(): Observable<boolean> {
    return this._isMaximized;
  }

  public toggleMaximized() {
    const currentWindow = remote.getCurrentWindow();
    if (currentWindow.isMaximized()) {
      currentWindow.unmaximize();
    } else {
      currentWindow.maximize();
    }
    this._isMaximized.next(currentWindow.isMaximized());
  }


  public minimize() {
    remote.getCurrentWindow().minimize();
  }

  public close() {
    remote.getCurrentWindow().close();
  }
}
