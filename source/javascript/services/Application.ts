module Application {

    export class TechnicConfig {
        public isSolderPack: boolean = true;

        public createForgeZip: boolean = false;
        public forgeVersion: string = "0";

        public checkPermissions: boolean = false;
        public isPublicPack: boolean = true;
    }

    export class FtbConfig {
        public isPublicPack: boolean = true;
    }

    export class Modpack {
        public name: string;
        public inputDirectory: string = "";
        public outputDirectory: string = "";
        public clearOutputDirectory: boolean = true;
        public minecraftVersion: string = "1.9";
        public version: string = "1.0.0";
        public additionalFolders: Array<string> = [];
        public technic: TechnicConfig = new TechnicConfig();
        public ftb: FtbConfig = new FtbConfig();

        constructor() {
            this.name = "Unnamed modpack";
            console.log(this);
        }
    }



    export class Application {
        static $inject = ["$rootScope"];
        public modpacks: Array<Modpack> = [];
        public modpack: Modpack;
        constructor(protected $rootScope: angular.IRootScopeService) {
            $rootScope.$watch(function() {
                return this.modpack;
            }, this.saveModpackData, true)
        }

        protected saveModpackData():void {

        }
    }



    angular.module("ModpackHelper").service("application", Application);
}
