module BuildController {
    export class BuildController {
        static $inject = ["application", "$mdDialog", "goComm", "$rootScope", "$translatePartialLoader"];
        
        public mods: Array<Application.Mod> = [];
        public todos: Array<string> = [];
        public showDone: boolean;
        public state: string = "info";
        public total: number;
        public progressNumber: number = 0;
        public uploading: string = "";
        public uploadNumber: number = 0;
        constructor(protected application: Application.Application, protected $mdDialog: angular.material.IDialogService, protected goComm: GoCommService.GoCommService, protected $rootScope: angular.IRootScopeService, protected $translatePartialLoader:angular.translate.ITranslatePartialLoaderService) {
            var self = this;
            $translatePartialLoader.addPart("build");

            $rootScope.$on("mod-data-ready", function (event: angular.IAngularEvent, mod: Application.Mod) {
                self.addModData(mod)
            });

            $rootScope.$on("total-to-pack", function(event: angular.IAngularEvent, total: number) {
                console.log("Total " + total);
                self.total = total;
            });
            
            $rootScope.$on("packing-part", function(event: angular.IAngularEvent, todo: string) {
                self.todos.push(todo)
            });

            $rootScope.$on("done-packing-part", function (event:angular.IAngularEvent, doneTodo:string) {
                var i = self.todos.indexOf(doneTodo);
                if(i > -1) {
                    self.todos.splice(i, 1)
                }
                self.progressNumber++;
            });

            $rootScope.$on("starting-upload", function(event:angular.IAngularEvent, key: string) {
                self.uploading = key;
                console.log("Starting upload of " + key);
            });

            $rootScope.$on("finished-uploading", function(event:angular.IAngularEvent, key: string) {
                self.uploadNumber++;
                console.log("Finished upload of " + key);
            });

            $rootScope.$on("finished-all-uploading", function() {
                self.uploading = "";
                console.log("Finished all uploading");
            });

            $rootScope.$on("started-uploading-all", function() {
                console.log("Starting uploading");
            });

            this.startBuild(application.modpack);
        }
        
        public startBuild(modpack: Application.Modpack): void {

            this.goComm.send("gather-information", modpack);
        }
        
        public addModData(mod: Application.Mod): void {
            var m = Application.Mod.fromJson(mod);
            if(!m.mcversion) {
                m.mcversion = this.application.modpack.minecraftVersion;
            }
            m.$$isDone = m.isValid();
            this.mods.push(m);
        }

        public cancel(): void{
            this.$mdDialog.hide();
        }
        
        public build(): void {
            var shouldBuild = true;
            for(var i = 0; i < this.mods.length; i++) {
                var mod = this.mods[i];
                mod.$$isDone = mod.isValid();
                if(!mod.$$isDone) {
                    shouldBuild = false;
                }
            }
            if(shouldBuild) {
                this.goComm.send("build", {modpack: this.application.modpack, mods: this.mods});
                this.state = "building"
            }
        }
    }

    angular.module("ModpackHelper").controller("BuildController", BuildController);
}
