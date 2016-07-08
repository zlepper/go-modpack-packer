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

    export class BuildController {
        static $inject = ["application", "$mdDialog", "goComm", "$rootScope", "$translatePartialLoader", "$window", "$document"];

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

        constructor(protected application:Application.Application,
                    protected $mdDialog:angular.material.IDialogService,
                    protected goComm:GoCommService.GoCommService,
                    protected $rootScope:angular.IRootScopeService,
                    protected $translatePartialLoader:angular.translate.ITranslatePartialLoaderService,
                    protected $window:angular.IWindowService) {
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
                console.log("Total " + total);
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
                console.log("Starting upload of " + key);
            });

            $rootScope.$on("finished-uploading", function (event:angular.IAngularEvent, key:string) {
                self.uploadNumber++;
                console.log("Finished upload of " + key);
            });

            $rootScope.$on("finished-all-uploading", function () {
                self.uploading = "";
                console.log("Finished all uploading");
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
                });
            });

            $rootScope.$on("waiting-for-file-upload", function (event:angular.IAngularEvent, data:UploadWaiting) {
                console.log(data);
                self.uploadData = data;
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
                if (!mod.$$isDone) {
                    shouldBuild = false;
                }
            }
            if (shouldBuild) {
                this.goComm.send("build", {modpack: this.application.modpack, mods: this.mods});
                this.state = "building"
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
            range.selectNodeContents($event.target);
            selection.removeAllRanges();
            selection.addRange(range);

            document.execCommand('copy');
            selection.removeAllRanges();
        }
    }

    angular.module("ModpackHelper").controller("BuildController", BuildController);
}
