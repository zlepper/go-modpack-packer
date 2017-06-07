import {Injectable} from "@angular/core";
import {Observable, Subject} from "rxjs";
import "../operators/pauseableBuffered";
import {ObservableWebsocket} from "./ObservableWebsocket";

export interface IWebsocketMessage {
  action: string;
  data: any;
}

@Injectable()
/**
 * Handles websocket connections to the backed.
 */
export class WebSocketService {
  /**
   * Indicates if the websocket is actually connected
   */
  private connectedStream: Observable<boolean>;

  private connection: ObservableWebsocket;
  private messageStream: Observable<IWebsocketMessage>;

  private messagesToSendSubject: Subject<IWebsocketMessage>;

  public constructor() {
    this.messagesToSendSubject = new Subject<IWebsocketMessage>();

    this.connection = new ObservableWebsocket('ws://localhost:8084/ws');
    console.log('websocket created');

    this.connectedStream = this.connection.events
      .filter(e => e.type === 'open' || e.type === 'close')
      .map(_ => this.connection.isConnected());

    this.connectedStream.subscribe(connected => console.log('websocket connected', connected));

    // Buffer messages when we don't have a connection
    this.messagesToSendSubject.pauseableBuffered(this.connectedStream.map(b => !b))
      .subscribe(message => {
        console.info('sending message to websocket', message);
        this.connection.send(JSON.stringify(message));
      });

    // Automatically reconnect
    this.connection.events
      .filter<CloseEvent>(e => e.type === 'close')
      .filter(e => !e.wasClean)
      .subscribe(() => this.reconnect());

    this.messageStream = this.connection.events
      .filter(e => e.type === 'message')
      .map((e: MessageEvent) => e.data)
      .map(data => {
        try {
          return <IWebsocketMessage>JSON.parse(data);
        } catch (err) {
          console.error(err);
        }
        return {};
      });

    this.connection.events
      .filter(e => e.type === 'error')
      .subscribe(e => console.error(e));
  }

  /**
   * A stream of messages provided by the underlying websocket.
   * The messages has already been translated into an object.
   * @returns {Observable<any>}
   */
  get messages(): Observable<IWebsocketMessage> {
    return this.messageStream;
  }


  get connected(): Observable<boolean> {
    return this.connectedStream;
  }

  /**
   * Sends a message over the websocket
   * @param message - The message to send
   */
  public send(message: IWebsocketMessage): void {
    this.messagesToSendSubject.next(message);
  }

  /**
   * Closes the connection to the remote server
   */
  public close(): void {
    this.connection.closeConnection();
  }

  /**
   * Reconnectes to the remote server
   */
  public reconnect() {
    this.connection.createConnection();
  }
}
