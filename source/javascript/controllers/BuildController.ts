module BuildController {
    class ModInfo {
        file:string;
        name:string;
        id:string;
        version:string;
        minecraftVersion:string;
        description:string;
        author:string;
        url:string;
        progressKey:string;
        permissions:Application.UserPermission;
    }

    class UploadWaiting {
        public modpack:Application.Modpack;
        public infos:Array<ModInfo>;
    }

    class Query {
        public order:string = "name";
        public page:number = 1;
        public limit:number = 10;
    }

    class PermissionSearch {
        public modId: string;
        public isPublic: boolean;

        constructor(id: string, isPublic: boolean) {
            this.modId = id;
            this.isPublic = isPublic;
        }
    }

    export class BuildController {
        static $inject = ["application", "$mdDialog", "goComm", "$rootScope", "$translatePartialLoader", "$window", "$translate", "$mdToast"];

        public mods:Array<Application.Mod> = [];
        public todos:Array<string> = [];
        public showDone:boolean;
        public state:string = "info";
        public total:number;
        public progressNumber:number = 0;
        public uploading:string = "";
        public uploadNumber:number = 0;
        public solderNumber:number = 0;
        public readyToBuild:boolean = false;
        public uploadData:UploadWaiting = null;
        public query:Query = new Query();
        public currentlyCheckingPermissions: { [id: string]: Application.Mod} = {};
        public permissionsText: string;
        public solderDoing: string;

        constructor(protected application:Application.Application,
                    protected $mdDialog:angular.material.IDialogService,
                    protected goComm:GoCommService.GoCommService,
                    protected $rootScope:angular.IRootScopeService,
                    protected $translatePartialLoader:angular.translate.ITranslatePartialLoaderService,
                    protected $window:angular.IWindowService,
                    protected $translate: angular.translate.ITranslateService,
                    protected $toast: angular.material.IToastService) {
            var self = this;
            $translatePartialLoader.addPart("build");


            $rootScope.$on("mod-data-ready", function (event:angular.IAngularEvent, mod:Application.Mod) {
                self.$window.requestAnimationFrame(function () {
                    self.addModData(mod)
                })
            });

            $rootScope.$on("all-mod-files-scanned", function () {
                self.readyToBuild = true;
            });

            $rootScope.$on("total-to-pack", function (event:angular.IAngularEvent, total:number) {
                self.total = total;
            });

            $rootScope.$on("packing-part", function (event:angular.IAngularEvent, todo:string) {
                self.$window.requestAnimationFrame(function () {
                    self.todos.push(todo)
                })
            });

            $rootScope.$on("done-packing-part", function (event:angular.IAngularEvent, doneTodo:string) {
                self.$window.requestAnimationFrame(function () {
                    var i = self.todos.indexOf(doneTodo);
                    if (i > -1) {
                        self.todos.splice(i, 1)
                    }
                    self.progressNumber++;
                });
            });

            $rootScope.$on("starting-upload", function (event:angular.IAngularEvent, key:string) {
                self.uploading = key;
            });

            $rootScope.$on("finished-uploading", function (event:angular.IAngularEvent, key:string) {
                self.uploadNumber++;
            });

            $rootScope.$on("finished-all-uploading", function () {
                self.uploading = "";
            });

            $rootScope.$on("started-uploading-all", function () {
                console.log("Starting uploading");
            });

            $rootScope.$on("updating-solder", function (event:angular.IAngularEvent, key:string) {
                self.$window.requestAnimationFrame(function () {
                    self.todos.push(key);
                });
            });

            $rootScope.$on("done-updating-solder", function (event:angular.IAngularEvent, key:string) {
                self.$window.requestAnimationFrame(function () {
                    var i = self.todos.indexOf(key);
                    if (i > -1) {
                        self.todos.splice(i, 1)
                    }
                    self.solderNumber++;
                    if (self.todos.length === 0) {
                        self.solderDoing = "BUILD.SOLDER.DONE";
                    }
                });
            });

            $rootScope.$on("waiting-for-file-upload", function (event:angular.IAngularEvent, data:UploadWaiting) {
                self.uploadData = data;
            });

            $rootScope.$on("got-permission-data", function (event:angular.IAngularEvent, data:Application.UserPermission) {
                var mod = self.currentlyCheckingPermissions[data.modId];
                data.modId = undefined;
                delete self.currentlyCheckingPermissions[mod.modid];
                mod.userPermission = data;
            });

            $rootScope.$on("solder-currently-doing", function(event:angular.IAngularEvent, status:string) {
                self.solderDoing = status;
            });

            this.startBuild(application.modpack);
        }

        public startBuild(modpack:Application.Modpack):void {

            this.goComm.send("gather-information", modpack);
        }

        public addModData(mod:Application.Mod):void {
            var m = Application.Mod.fromJson(mod);
            if (!m.mcversion) {
                m.mcversion = this.application.modpack.minecraftVersion;
            }
            m.$$isDone = m.isValid();
            this.mods.push(m);
        }

        public cancel():void {
            this.$mdDialog.hide();
        }

        public build():void {
            var shouldBuild = true;
            for (var i = 0; i < this.mods.length; i++) {
                var mod = this.mods[i];
                mod.$$isDone = mod.isValid();
                if (!mod.$$isDone && !mod.skip) {
                    shouldBuild = false;
                }
            }
            if (shouldBuild) {
                this.goComm.send("build", {modpack: this.application.modpack, mods: this.mods.filter((m) => !m.skip)});
                this.state = "building"
            } else {
                this.$translate('BUILD.MOD.VALIDATION_FAILED').then((t) => {
                    this.$toast.showSimple(t)
                });
            }
        }

        public continueRunning():void {
            var uploadData = this.uploadData;
            this.goComm.send("continue-running", uploadData);
            this.uploadData = null;
        }

        public copyInput($event: JQueryEventObject):void {
            var selection = this.$window.getSelection();
            var range = document.createRange();
            // $event.currentTarget does not exist on the object, so just ignore how this line works,
            // because honest truth, i have no clue what so ever. 
            range.selectNodeContents(<any>$event.currentTarget);
            selection.removeAllRanges();
            selection.addRange(range);

            document.execCommand('copy');
            selection.removeAllRanges();
        }

        public checkDbForPermissions(mod: Application.Mod) {
            // Don't send request for something for which we are currently checking.
            if(this.currentlyCheckingPermissions[mod.modid]) return;
            this.currentlyCheckingPermissions[mod.modid] = mod;
            this.goComm.send("check-permission-store", new PermissionSearch(mod.modid, this.application.modpack.technic.isPublicPack));
        }
    }



    angular.module("ModpackHelper").controller("BuildController", BuildController);
}
