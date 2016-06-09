module TechnicController {
    export class TechnicController {
        static $inject = ["application", "$translatePartialLoader", "$mdDialog", "$mdMedia", "forge", "goComm", "$rootScope"];

        public buckets: Array<string> = [];

        public ftpPattern: RegExp = /^[\w\.]+:\d+$/;

        constructor(protected application: Application.Application,
                    protected $translatePartialLoader: angular.translate.ITranslatePartialLoaderService,
                    protected $mdDialog: angular.material.IDialogService,
                    protected $mdMedia: angular.material.IMedia,
                    protected forge: ForgeVersion.ForgeVersionService,
                    protected goComm: GoCommService.GoCommService,
                    protected $rootScope: angular.IRootScopeService) {
            $translatePartialLoader.addPart("technic");
            var controller = this;
            $rootScope.$on("found-aws-buckets", function(event: angular.IAngularEvent, buckets: Array<string>) {
                controller.buckets = buckets;
            });
        }

        public build(ev: MouseEvent): void {
            if(this.application.modpack.isValid() && this.application.modpack.isValidTechnic()) {
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
                    this.$mdDialog.alert() // TODO Get translations
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

        public getAwsBuckets():void {
            console.log("GEt aws buckets");
            this.goComm.send("get-aws-buckets", this.application.modpack)
        }
    }

    angular.module("ModpackHelper").controller("TechnicController", TechnicController);
}
