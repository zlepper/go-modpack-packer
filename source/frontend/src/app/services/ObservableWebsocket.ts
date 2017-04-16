import {Observable, Subject} from "rxjs";
/**
 * A websocket wrapper that will handle automatic reconnecting without
 * losing observables
 */
export class ObservableWebsocket {
  private websocketEvents: Subject<Event>;
  private connection: WebSocket;
  private readonly url: string;

  public constructor(url: string) {
    this.websocketEvents = new Subject<Event>();

    this.url = url;

    this.createConnection();
  }

  public createConnection() {
    this.connection = new WebSocket(this.url);

    this.connection.onmessage = e => this.websocketEvents.next(e);
    this.connection.onopen = e => this.websocketEvents.next(e);
    this.connection.onerror = e => this.websocketEvents.next(e);
    this.connection.onclose = e => this.websocketEvents.next(e);
  }

  public closeConnection() {
    if (this.connection) {
      this.connection.close();
      // Prevent more events from being emitted
      this.connection.onmessage = null;
      this.connection.onopen = null;
      this.connection.onerror = null;
      this.connection.onclose = null;
      this.connection = null;
    }
  }

  public isConnected(): boolean {
    return this.connection.readyState === WebSocket.OPEN;
  }

  public get events(): Observable<Event> {
    return this.websocketEvents;
  }

  public send(message: string): void {
    this.connection.send(message);
  }
}
