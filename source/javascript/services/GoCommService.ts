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
        static $inject = ["$websocket", "$rootScope", "$timeout", "$interval"];

        private dataStream: IWebsocket;
        private ready: boolean;
        private events: Array<IWebsocketOnMessageEvent> = [];
        constructor(private $websocket: IWebsocketService, protected $rootScope: angular.IRootScopeService, protected $timeout: angular.ITimeoutService, protected $interval: angular.IIntervalService) {
            var t = this;
            this.dataStream = $websocket("ws://localhost:8084/ws");
            this.dataStream.onOpen(function() {
                t.ready = true;
            });

            this.$interval(function() {
                if(t.events.length > 0) {
                    var data = t.events.shift();
                    var message:MessageData = JSON.parse(data.data);
                    // Special logging trick
                    if (message.action === "log") {
                        return console.log(message.data);
                    }
                    $rootScope.$emit(message.action, message.data);
                }
            }, 5);

            this.dataStream.onMessage(function(data: IWebsocketOnMessageEvent) {
                t.events.push(data);
            });
        }

        public send(action: string, data: any) {
            if(this.ready) {
                var md = new MessageData();
                md.action = action;
                md.data = data;
                this.dataStream.send(angular.toJson(md));
            } else {
                var t = this;
                // Keep retrying in case the socket is not yet ready
                this.$timeout(function() {
                    t.send(action, data);
                }, 50);
            }
        }


    }

    angular.module("ModpackHelper").service("goComm", GoCommService);
}
