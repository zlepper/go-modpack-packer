module MainController {
    import Modpack = Application.Modpack;
    var remote = require("electron").remote;
    export class MainController {
        public static $inject = ["application", "$state", "electron"];
        public isMaximized: boolean = false;

        constructor(protected application:Application.Application, 
                    protected $state:angular.ui.IStateService,
                    protected electron: ElectronService.ElectronService) {
                
        }

        public createNewModpack():void {
            var modpack = new Application.Modpack();

            this.application.modpacks.push(modpack);

            this.application.modpack = modpack;

            this.$state.go("modpack");
            
            this.isMaximized = remote.getCurrentWindow().isMaximized();
        }
        
        public restart():void {
            this.electron.send("restart", null);
        }
        
        public close():void {
            remote.getCurrentWindow().close();
        }
        
        public toggleMaximized():void {
            if (this.isMaximized) {
                remote.getCurrentWindow().unmaximize();
            } else {
                remote.getCurrentWindow().maximize()
            }
            this.isMaximized = remote.getCurrentWindow().isMaximized();
        }

        public minimize():void {
            remote.getCurrentWindow().minimize();
        }
    }

    angular.module("ModpackHelper").controller("MainController", MainController);
}
