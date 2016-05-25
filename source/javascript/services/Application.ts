module Application {

    export class TechnicConfig {
        public isSolderPack:boolean = true;

        public createForgeZip:boolean = false;
        public forgeVersion:string = "0";

        public checkPermissions:boolean = false;
        public isPublicPack:boolean = true;
    }

    export class FtbConfig {
        public isPublicPack:boolean = true;
    }

    export class Modpack {
        public name:string;
        public inputDirectory:string = "";
        public outputDirectory:string = "";
        public clearOutputDirectory:boolean = true;
        public minecraftVersion:string = "1.9";
        public version:string = "1.0.0";
        public additionalFolders:{[folder:string]:boolean } = {};
        public technic:TechnicConfig = new TechnicConfig();
        public ftb:FtbConfig = new FtbConfig();

        constructor() {
            this.name = "Unnamed modpack";
            console.log(this);
        }

        public static fromJson(data:Modpack):Modpack {
            var modpack = new Modpack();
            modpack.name = data.name;
            modpack.inputDirectory = data.inputDirectory;
            modpack.outputDirectory = data.outputDirectory;
            modpack.clearOutputDirectory = data.clearOutputDirectory;
            modpack.minecraftVersion = data.minecraftVersion;
            modpack.version = data.version;
            modpack.additionalFolders = data.additionalFolders;
            modpack.technic = data.technic;
            modpack.ftb = data.ftb;
            return modpack;
        }
    }


    export class Application {
        static $inject = ["$rootScope", "goComm", "$state"];
        public modpacks:Array<Modpack> = [];
        public modpack:Modpack;

        constructor(protected $rootScope:angular.IRootScopeService, protected goComm:GoCommService.GoCommService, protected $state: angular.ui.IStateService) {
            var self = this;
            goComm.send("load-modpacks", {});
            $rootScope.$watch(function () {
                return self.modpacks;
            }, function () {
                self.saveModpackData()
            }, true);
            $rootScope.$on("data-loaded", (event:angular.IAngularEvent, modpacks:Array<Modpack>) => {
                console.log("Data loaded");

                modpacks.forEach((modpack:Modpack) => {
                    self.modpacks.push(Modpack.fromJson(modpack))
                });
                if (self.modpacks.length) {
                    self.modpack = self.modpacks[0];
                    self.$state.go("modpack");
                }
            });
        }

        protected saveModpackData():void {
            this.goComm.send("save-modpacks", this.modpacks);
        }
    }


    angular.module("ModpackHelper").service("application", Application);
}
