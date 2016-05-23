module Electron {

    import IpcRendererEvent = Electron.IpcRendererEvent;
    export class ElectronService {
        static $inject = ["$timeout"];
        private ipc = require("electron").ipcRenderer;
        private electron_host = "ELECTRON_HOST";

        constructor(public $timeout:angular.ITimeoutService) {

        }

        public send(data:string) {
            this.ipc.send(data);
        }

        // Totally not a hack to send things back into the angular event loop
        public on(channel:string, cb:Electron.IpcRendererEventListener) {
            this.ipc.on(channel, (event:IpcRendererEvent, ...args:any[]) => {
                this.$timeout(function() {
                    cb(event, args);
                }, 0);
            });
        }
    }

    angular.module("ModpackHelper").service("electron", ElectronService);
}
