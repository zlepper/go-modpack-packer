module TechnicController {
    export class TechnicController {
        static $inject = ["application", "$translatePartialLoader", "$mdDialog", "$mdMedia", "forge"];
        
        constructor(protected application: Application.Application, protected $translatePartialLoader: angular.translate.ITranslatePartialLoaderService, protected $mdDialog: angular.material.IDialogService, protected $mdMedia: angular.material.IMedia, protected forge: ForgeVersion.ForgeVersionService) {
            $translatePartialLoader.addPart("technic");
        }

        public build(ev: MouseEvent): void {
            if(this.application.modpack.isValid()) {
                var useFullscreen = this.$mdMedia("sm") || this.$mdMedia("xs");
                this.$mdDialog.show({
                    controller: "BuildController",
                    controllerAs: "build",
                    templateUrl: "parts/buildprogress.html",
                    clickOutsideToClose: false,
                    fullscreen: useFullscreen,
                    targetEvent: ev,
                    ariaLabel: "Build progress dialog"
                });
            } else {
                this.$mdDialog.show(
                    this.$mdDialog.alert()
                        .clickOutsideToClose(true)
                        .title("Missing info")
                        .textContent("The modpack is missing some info before it can be build.")
                        .ariaLabel("Missing information")
                        .ok("I will go fix my mistakes")
                        .targetEvent(ev)
                )
            }
        }
        
        public filterByMcVersion(input: ForgeVersion.ForgeVersion, minecraftVersion: string) {
            return input.minecraftVersion === minecraftVersion;
        }
    }

    angular.module("ModpackHelper").controller("TechnicController", TechnicController);
}
