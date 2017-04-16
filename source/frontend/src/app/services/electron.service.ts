import {ApplicationRef, Injectable} from "@angular/core";
import {BehaviorSubject, Observable, Subject} from "rxjs";

declare var electron;
let remote;
if(typeof electron === 'undefined') {
    remote = {
      getCurrentWindow() {
        return {
          isMaximized() {
            return false;
          }
        }
      }
    };
} else {
  remote = electron.remote;
}

export interface IPCMessage {
  key: string;
  data: any;
}

@Injectable()
export class ElectronService {
  private _isMaximized: Subject<boolean>;
  private ipc = electron.ipcRenderer;

  constructor(protected applicationRef: ApplicationRef) {
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

  public send(key: string, data: any) {
    this.ipc.send(key, data);
  }

  public on(channel: string, cb): void {
    this.ipc.on(channel, (event, ...args) => {
      cb(...args);
    });
  }

  public once(channel: string, cb): void {
    this.ipc.once(channel, (event, ...args) => {
      cb(...args);
    });
  }
}
