module GoCommService {
    class MessageData {
        public action: string;
        public data: any;
    }

    interface IWebsocketOnMessageEvent {
        data: string;
    }

    interface OnMessageCallback {
        (data: IWebsocketOnMessageEvent): void;
    }

    interface OnOpenCallback {
        (): void;
    }

    interface IWebsocket {
        onMessage(cb: OnMessageCallback): void;
        send(data: string): void;
        onOpen(db: OnOpenCallback): void;
    }

    interface IWebsocketService {
        (url: string): IWebsocket;
    }

    export class GoCommService {
        static $inject = ["$websocket", "$rootScope", "$timeout"];

        private dataStream: IWebsocket;
        private ready: boolean;
        constructor(private $websocket: IWebsocketService, protected $rootScope: angular.IRootScopeService, protected $timeout: angular.ITimeoutService) {
            var t = this;
            this.dataStream = $websocket("ws://localhost:8084/ws");
            this.dataStream.onOpen(function() {
                console.log("Websocket is ready");
                t.ready = true;
            });
            this.dataStream.onMessage(function(data: IWebsocketOnMessageEvent) {
                var message: MessageData = JSON.parse(data.data);
                // Special logging trick
                if(message.action === "log") {
                    return console.log(message.data);
                }
                $rootScope.$emit(message.action, message.data);
            });
        }

        public send(action: string, data: any) {
            if(this.ready) {
                console.log("Sending to websocket");
                var md = new MessageData();
                md.action = action;
                md.data = data;
                this.dataStream.send(JSON.stringify(md));
            } else {
                var t = this;
                // Keep retrying in case the socket is not yet ready
                this.$timeout(function() {
                    console.log("Waiting");
                    t.send(action, data);
                }, 50);
            }
        }


    }

    angular.module("ModpackHelper").service("goComm", GoCommService);
}
