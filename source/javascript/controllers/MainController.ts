module MainController {
    import Modpack = Application.Modpack;
    export class MainController {
        public $inject = ["application", "$state"];

        constructor(protected application:Application.Application, protected $state:angular.ui.IStateService) {

        }

        public createNewModpack():void {
            var modpack = new Application.Modpack();

            this.application.modpacks.push(modpack);

            this.application.modpack = modpack;

            this.$state.go("modpack");
        }
    }

    angular.module("ModpackHelper").controller("MainController", MainController);
}
