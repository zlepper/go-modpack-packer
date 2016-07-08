module Application {

    export class AWSConfig {
        public accessKey:string = "";
        public secretKey:string = "";
        public region:string = "us-east-1";
        public bucket:string = "";
    }

    export class FtpConfig {
        public url:string = "";
        public username:string = "";
        public password:string = "";
    }

    export class UploadConfig {
        public type:string = "none";
        public aws:AWSConfig = new AWSConfig();
        public ftp:FtpConfig = new FtpConfig();
    }


    export class TechnicConfig {
        public isSolderPack:number = 1;

        public createForgeZip:boolean = false;
        public forgeVersion:ForgeVersion.ForgeVersion;

        public checkPermissions:boolean = false;
        public isPublicPack:boolean = true;

        public memory:number = 0;
        public java:string = "1.8";

        public upload:UploadConfig = new UploadConfig();

        public repackAllMods: boolean = false;
    }

    export class SolderInfo {
        public use:boolean = false;
        public url:string = "";
        public username:string = "";
        public password:string = "";
    }

    export class FtbConfig {
        public isPublicPack:boolean = true;
    }

    export class Folder {
        public name:string;
        public include:boolean;
    }

    export class Modpack {
        public name:string;
        public inputDirectory:string = "";
        public outputDirectory:string = "";
        public clearOutputDirectory:boolean = true;
        public minecraftVersion:string = "1.9";
        public version:string = "1.0.0";
        public additionalFolders:Array<Folder> = [];
        public technic:TechnicConfig = new TechnicConfig();
        public ftb:FtbConfig = new FtbConfig();
        public solder:SolderInfo = new SolderInfo();

        constructor() {
            this.name = "Unnamed modpack";
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
            modpack.solder = data.solder;
            return modpack;
        }

        public isValid():boolean {
            if (!this.name) return false;
            if (!this.inputDirectory) return false;
            if (!this.outputDirectory) return false;
            if (!this.minecraftVersion) return false;
            if (!this.version) return false;
            return true;
        }

        public isValidTechnic():boolean {
            var t = this.technic;
            console.log(t.upload.type);
            switch (t.upload.type) {
                case "s3":
                    var aws = t.upload.aws;
                    if(!aws.accessKey || !aws.bucket || !aws.region || !aws.secretKey) {
                        return false;
                    }
                    break;
                case "ftp":
                    var ftp = t.upload.ftp;
                    // Don't validate password existing, because it's possible to connect to ftp without a password
                    // even though that is a very bad idea. Security and all that.
                    if(!ftp.url || !ftp.username) {
                        return false;
                    }
                    break;
                case "none":
                default:
                    break;
            }
            if(t.isSolderPack == 1) {
                var solder = this.solder;
                if(solder.use) {
                    if (!solder.url ||!solder.password || !solder.username) {
                        return false;
                    }
                }
            }
            return true;
        }
    }

    export class Mod {
        public modid:string;
        public name:string;
        public description:string;
        public version:string;
        public mcversion:string;
        public url:string;
        public authors:string;
        public credits:string;
        public filename:string;
        public md5: string;
        // Naming is totally a hack to make sure the value does not get send to the server
        public $$isDone:boolean;
        public isOnSolder: boolean;

        public static fromJson(data:Mod):Mod {
            var m = new Mod();
            m.modid = data.modid;
            m.name = data.name;
            m.description = data.description;
            m.version = data.version;
            m.mcversion = data.mcversion;
            m.url = data.url;
            m.authors = data.authors;
            m.credits = data.credits;
            m.filename = data.filename;
            m.md5 = data.md5;
            m.isOnSolder = data.isOnSolder;
            return m;
        }

        public isValid():boolean {
            if (!this.modid) return false;
            if (!this.name) return false;
            if (!this.version) return false;
            if (!this.mcversion) return false;
            if (this.authors.length < 1) return false;

            return this.isAdvancedValid();
        }

        private isAdvancedValid(): boolean {
            if(this.modid.toLowerCase().indexOf("example") > -1) {
                return false;
            }
            if(this.name.toLowerCase().indexOf("example") > -1) {
                return false;
            }
            if(this.version.toLowerCase().indexOf("example") > -1) {
                return false;
            }
            if(this.name.indexOf("${") > -1) {
                return false;
            }
            if(this.version.indexOf("${") > -1) {
                return false;
            }
            if(this.mcversion.indexOf("${") > -1) {
                return false;
            }
            if(this.modid.indexOf("${") > -1) {
                return false;
            }
            if(this.version.toLowerCase().indexOf("@version@") > -1) {
                return false;
            }
            return true;
        }
    }

    var electron = require("electron");

    export class Application {
        static $inject = ["$rootScope", "goComm", "$state", "$mdToast", "$translate", "$timeout", "languageService"];
        public modpacks:Array<Modpack> = [];
        public modpack:Modpack;
        public waitingForStoredData: boolean = false;
        public updateReady: boolean = false;
        constructor(protected $rootScope:angular.IRootScopeService,
                    protected goComm:GoCommService.GoCommService,
                    protected $state:angular.ui.IStateService,
                    protected $mdToast: angular.material.IToastService,
                    protected $translate: angular.translate.ITranslateService,
                    protected $timeout: angular.ITimeoutService,
                    protected languageService: LanguageService.LanguageService) {
            var self = this;
            goComm.send("load-modpacks", {});
            this.waitingForStoredData = true;
            $rootScope.$watch(function () {
                return self.modpacks;
            }, function () {
                self.saveModpackData()
            }, true);
            
            $timeout(function wait() {
                if(self.waitingForStoredData) {
                    goComm.send("load-modpacks", {});
                    $timeout(wait, 5000)
                }
            }, 5000);
            
            $rootScope.$on("data-loaded", (event:angular.IAngularEvent, modpacks:Array<Modpack>) => {
                self.waitingForStoredData = false;
                console.log("Got stored data");
                modpacks.forEach((modpack:Modpack) => {
                    self.modpacks.push(Modpack.fromJson(modpack))
                });
                if (self.modpacks.length) {
                    self.modpack = self.modpacks[0];
                    if (self.$state.is("home")) {
                        self.$state.go("modpack");
                    }
                }
            });

            $rootScope.$on("error", (event:angular.IAngularEvent, err: string) => {
                console.log(err);
                var self = this;
                this.$translate(err).then(function(translation: string) {
                    self.$mdToast.showSimple(translation);
                });
            });
            
            electron.ipcRenderer.on("update-info", (event: Electron.IpcRendererEvent, message: string) => {
                self.$translate(message).then(function(translated) {
                    $mdToast.showSimple(translated);
                });
                
                if (message === "UPDATE.DOWNLOADED") {
                    self.updateReady = true;
                }
            });
        }

        protected saveModpackData():void {
            this.goComm.send("save-modpacks", this.modpacks);
        }
    }


    angular.module("ModpackHelper").service("application", Application);


    electron.ipcRenderer.on("update-error", (event: Electron.IpcRendererEvent, message: any) => {
        console.error(message);
    })
}
