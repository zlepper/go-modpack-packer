module MainController {
    import Modpack = Application.Modpack;
    var remote = require("electron").remote;
    export class MainController {
        public static $inject = ["application", "$state", "electron", "$translate", "$mdToast"];
        public isMaximized: boolean = false;

        constructor(protected application:Application.Application, 
                    protected $state:angular.ui.IStateService,
                    protected electron: ElectronService.ElectronService,
                    protected $translate: angular.translate.ITranslateService,
                    protected $mdToast: angular.material.IToastService) {

            this.electron.on('error', (_: any, err: any) => {
                console.log(err[0]);
                this.$mdToast.show(
                    this.$mdToast.simple().textContent("Background process crashed. Please report an issue on the bugtracker.")
                        .hideDelay(0)
                );
            });
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

        public selectModpack() {
            var modpack = this.application.modpack;
            if(modpack && modpack.isNew) {
                modpack.isNew = false;
                this.$translate("MODPACK.UNNAMED").then(t => {
                    modpack.name = t;
                });
                this.$state.go('modpack');

                this.application.addNewModpack();

            }
        }
    }

    angular.module("ModpackHelper").controller("MainController", MainController);
}
