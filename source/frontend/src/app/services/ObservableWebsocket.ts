import {Observable, Subject} from "rxjs";
import {fromEvent} from "rxjs/observable/fromEvent";
import {Subscription} from "rxjs/Subscription";

/**
 * A websocket wrapper that will handle automatic reconnecting without
 * losing observables
 */
export class ObservableWebsocket {
  private websocketEvents: Subject<Event>;
  private connection?: WebSocket;
  private readonly url: string;
  private websocketSubscription: Subscription;

  public constructor(url: string) {
    this.websocketEvents = new Subject<Event>();

    this.url = url;

    this.createConnection();
  }

  public createConnection() {
    this.connection = new WebSocket(this.url);
    this.websocketSubscription = fromEvent(this.connection, 'message')
      .merge(fromEvent(this.connection, 'open'))
      .merge(fromEvent(this.connection, 'error'))
      .merge(fromEvent(this.connection, 'close'))
      .subscribe((e: Event) => {
        console.debug(e);
        return this.websocketEvents.next(e);
      });
  }

  public closeConnection() {
    console.trace('closing websocket connection');
    if (this.connection && !this.websocketSubscription.closed) {
      this.connection.close();
      this.websocketSubscription.unsubscribe();
      this.connection = undefined;
    }
  }

  public isConnected(): boolean {

    return !!this.connection && this.connection.readyState === WebSocket.OPEN;
  }

  public get events(): Observable<Event> {
    return this.websocketEvents;
  }

  public send(message: string): void {
    if (this.connection) {
      this.connection.send(message);
    }
  }
}
