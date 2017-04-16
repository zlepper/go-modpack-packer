import {Injectable} from "@angular/core";
import {MdSnackBar} from "@angular/material";
import {WebSocketService} from "app/services/websocket.service";
import {Observable} from "rxjs";

@Injectable()
export class BackendCommunicationService {

  constructor(protected websocketService: WebSocketService, protected snackBar: MdSnackBar) {
    this.websocketService.messages.filter(message => message.action === 'log')
      .subscribe(message => console.log(message.data));
    this.websocketService.messages.filter(message => message.action === 'notification')
      .subscribe(message => snackBar.open(message.data, null, {duration: 5000}));
  }

  public send(action: string, data: any) {
    this.websocketService.send({action, data})
  }

  public getMessages<T>(type: string): Observable<T> {
    return this.websocketService.messages
      .filter(message => message.action === type)
      .map(message => <T>message.data);
  }

}
