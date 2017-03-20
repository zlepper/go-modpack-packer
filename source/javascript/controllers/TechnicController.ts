module TechnicController {
    export class TechnicController {
        static $inject = ["application", "$translatePartialLoader", "$mdDialog", "$mdMedia", "forge", "goComm", "$rootScope", "$translate", "$mdToast"];

        public buckets: Array<string> = [];

        public ftpPattern: RegExp = /^[\d\w.]+(?::(?:\d){1,5})$/;

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
            var validation = this.application.modpack.isValid();
            if(!validation) {
                validation = this.application.modpack.isValidTechnic();
            }
            if(!validation) {
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
                this.$translate(validation).then(validation => {
                    this.$translate("TECHNIC.ERRORS.MISSING.INFO.TITLE").then(title => {
                        this.$translate("TECHNIC.ERRORS.MISSING.INFO.BODY").then(body => {
                            this.$translate("TECHNIC.ERRORS.MISSING.INFO.OK_BUTTON").then(ok => {
                                this.$mdDialog.show(
                                    this.$mdDialog.alert()
                                        .clickOutsideToClose(true)
                                        .title(title)
                                        .textContent(body + "\n" + validation)
                                        .ariaLabel(title)
                                        .ok(ok)
                                        .targetEvent(ev)
                                )
                            });
                        });
                    });
                });
            }
        }
        public testFtp(): void {
            this.goComm.send("test-ftp", this.application.modpack.technic.upload.ftp)
        }
        public testSolder(): void {
            let re = /.*\/api\/?$/;
            if(re.test(this.application.modpack.solder.url)) {
                this.$translate("TECHNIC.ERRORS.SOLDER_URL_SHOULD_NOT_CONTAIN_API").then(t => {
                    this.$mdToast.showSimple(t)
                });
            } else {
                this.goComm.send("test-solder", this.application.modpack.solder)
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
