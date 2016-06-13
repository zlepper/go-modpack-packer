module TechnicController {
    export class TechnicController {
        static $inject = ["application", "$translatePartialLoader", "$mdDialog", "$mdMedia", "forge", "goComm", "$rootScope", "$translate", "$mdToast"];

        public buckets: Array<string> = [];

        public ftpPattern: RegExp = /^[\w\.]+:\d+$/;

        constructor(protected application: Application.Application,
                    protected $translatePartialLoader: angular.translate.ITranslatePartialLoaderService,
                    protected $mdDialog: angular.material.IDialogService,
                    protected $mdMedia: angular.material.IMedia,
                    protected forge: ForgeVersion.ForgeVersionService,
                    protected goComm: GoCommService.GoCommService,
                    protected $rootScope: angular.IRootScopeService,
                    protected $translate: angular.translate.ITranslateService,
                    protected $mdToast: angular.material.IToastService) {
            $translatePartialLoader.addPart("technic");
            var controller = this;
            $rootScope.$on("found-aws-buckets", function(event: angular.IAngularEvent, buckets: Array<string>) {
                controller.buckets = buckets;
            });
            $rootScope.$on("solder-test", function(event: angular.IAngularEvent, result: string) {
                $translate(result).then(function(translated) {
                    $mdToast.showSimple(translated);
                });
            });

        }

        public build(ev: MouseEvent): void {
            if(this.application.modpack.isValid() && this.application.modpack.isValidTechnic()) {
                this.$mdDialog.show({
                    controller: "BuildController",
                    controllerAs: "build",
                    templateUrl: "parts/buildprogress.html",
                    clickOutsideToClose: false,
                    fullscreen: true,
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
        public testFtp(): void {
            this.goComm.send("test-ftp", this.application.modpack.technic.upload.ftp)
        }
        public testSolder(): void {
            this.goComm.send("test-solder", this.application.modpack.solder)
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
