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
        public path:string = "";
    }

    export class UploadConfig {
        public type:string = "none";
        public aws:AWSConfig = new AWSConfig();
        public ftp:FtpConfig = new FtpConfig();
    }


    export class TechnicConfig {
        public isSolderPack:boolean = true;

        public createForgeZip:boolean = false;
        public forgeVersion:ForgeVersion.ForgeVersion;

        public checkPermissions:boolean = false;
        public isPublicPack:boolean = true;

        public memory:number = 0;
        public java:string = "1.8";

        public upload:UploadConfig = new UploadConfig();

        public repackAllMods:boolean = false;
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
        public isNew:boolean = false;
        public $$hash: number;

        constructor() {
            this.$$hash = Math.floor(Math.random() * 100000);
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
            modpack.isNew = data.isNew;
            modpack.$$hash = Math.floor(Math.random() * 100000);
            return modpack;
        }

        /**
         * Checks if the basic modpack info is valid.
         *
         * Returns an empty string if the basic modpack info is valid, otherwise an error message key.
         * @returns {string}
         */
        public isValid():string {
            if (!this.name) return "MODPACK.ERRORS.MISSING.NAME";
            if (!this.inputDirectory) return "MODPACK.ERRORS.MISSING.INPUT_DIRECTORY";
            if (!this.outputDirectory) return "MODPACK.ERRORS.MISSING.OUTPUT_DIRECTORY";
            if (!this.minecraftVersion) return "MODPACK.ERRORS.MISSING.MINECRAFT_VERSION";
            if (!this.version) return "MODPACK.ERRORS.MISSING.VERSION";
            return "";
        }

        public isValidTechnic():string {
            var t = this.technic;
            console.log(t.upload.type);
            switch (t.upload.type) {
                case "s3":
                    var aws = t.upload.aws;
                    if (!aws.accessKey) {
                        return "TECHNIC.ERRORS.MISSING.AWS.ACCESS_KEY";
                    }
                    if (!aws.bucket) {
                        return "TECHNIC.ERRORS.MISSING.AWS.BUCKET";
                    }
                    if (!aws.region) {
                        return "TECHNIC.ERRORS.MISSING.AWS.REGION";
                    }
                    if (!aws.secretKey) {
                        return "TECHNIC.ERRORS.MISSING.AWS.SECRET_KEY";
                    }
                    break;
                case "ftp":
                    var ftp = t.upload.ftp;
                    // Don't validate password existing, because it's possible to connect to ftp without a password
                    // even though that is a very bad idea. Security and all that.
                    if (!ftp.url) {
                        return "TECHNIC.ERRORS.MISSING.FTP.URL";
                    }
                    if (!ftp.username) {
                        return "TECHNIC.ERRORS.MISSING.FTP.USERNAME";
                    }
                    break;
                case "none":
                default:
                    break;
            }
            if (t.isSolderPack) {
                var solder = this.solder;
                if (solder.use) {
                    if (!solder.url) {
                        return "TECHNIC.ERRORS.MISSING.SOLDER.URL";
                    }
                    let re = /.*\/api\/?$/;
                    if(re.test(solder.url)) {
                        return "TECHNIC.ERRORS.SOLDER_URL_SHOULD_NOT_CONTAIN_API"
                    }
                    if (!solder.password) {
                        return "TECHNIC.ERRORS.MISSING.SOLDER.PASSWORD";
                    }
                    if (!solder.username) {
                        return "TECHNIC.ERRORS.MISSING.SOLDER.USERNAME";
                    }
                }
            }
            return "";
        }
    }

    export class UserPermission {
        public licenseLink:string;
        public modLink:string;
        public permissionLink:string;
        public policy:string;
        public modId:string;
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
        public md5:string;
        // Naming is totally a hack to make sure the value does not get send to the server
        public $$isDone:boolean;
        public isOnSolder:boolean;
        public userPermission:UserPermission;
        public skip:boolean;

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
            m.userPermission = data.userPermission;
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

        private isAdvancedValid():boolean {
            if (this.modid.toLowerCase().indexOf("example") > -1) {
                return false;
            }
            if (this.name.toLowerCase().indexOf("example") > -1) {
                return false;
            }
            if (this.version.toLowerCase().indexOf("example") > -1) {
                return false;
            }
            if (this.name.indexOf("${") > -1) {
                return false;
            }
            if (this.version.indexOf("${") > -1) {
                return false;
            }
            if (this.mcversion.indexOf("${") > -1) {
                return false;
            }
            if (this.modid.indexOf("${") > -1) {
                return false;
            }
            if (this.version.toLowerCase().indexOf("@version@") > -1) {
                return false;
            }
            if (this.userPermission) {
                if (this.userPermission.policy !== "Open") {
                    if (!this.userPermission.licenseLink) return false;
                    if (!this.userPermission.modLink) return false;
                    if (!this.userPermission.permissionLink) return false;
                }
            }

            return true;
        }
    }

    var electron = require("electron");

    export class Application {
        static $inject = ["$rootScope", "goComm", "$state", "$mdToast", "$translate", "$timeout", "languageService"];
        public modpacks:Array<Modpack> = [];
        public modpack:Modpack;
        public waitingForStoredData:boolean = false;
        public updateReady:boolean = false;

        constructor(protected $rootScope:angular.IRootScopeService,
                    protected goComm:GoCommService.GoCommService,
                    protected $state:angular.ui.IStateService,
                    protected $mdToast:angular.material.IToastService,
                    protected $translate:angular.translate.ITranslateService,
                    protected $timeout:angular.ITimeoutService,
                    protected languageService:LanguageService.LanguageService) {
            var self = this;
            goComm.send("load-modpacks", {});
            this.waitingForStoredData = true;
            $rootScope.$watch(function () {
                return self.modpacks;
            }, function () {
                if (self.modpacks && self.modpacks.length > 0) {
                    self.saveModpackData()
                }
            }, true);

            $timeout(function wait() {
                if (self.waitingForStoredData) {
                    goComm.send("load-modpacks", {});
                    $timeout(wait, 5000)
                }
            }, 5000);

            $rootScope.$on("data-loaded", (event:angular.IAngularEvent, modpacks:Array<Modpack>) => {
                if(!self.waitingForStoredData) return;

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
                    // If the last modpack in the list is already a "new" pack, then we shouldn't append another.
                    if(self.modpacks[self.modpacks.length - 1].isNew) {
                         return;
                    }
                }
                console.log(self.modpacks);
                
             
                this.addNewModpack();
                // Select the newly created modpack
                this.modpack = this.modpacks[0];
                this.$translate("MODPACK.UNNAMED").then(t => {
                    this.modpack.name = t;
                });
                this.modpack.isNew = false;

                // Add another modpack the user can select to create a new one. 
                this.addNewModpack();
            });

            $rootScope.$on("error", (event:angular.IAngularEvent, err:string) => {
                console.log(err);
                var self = this;
                this.$translate(err).then(function (translation:string) {
                    self.$mdToast.showSimple(translation);
                });
            });

            electron.ipcRenderer.on("update-info", (event:Electron.IpcRendererEvent, message:string) => {
                self.$translate(message).then(function (translated) {
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

        public addNewModpack():void {
            var modpack = new Modpack();
            this.$translate("MODPACK.NEW").then(t => {
                modpack.name = t;
            });
            modpack.isNew = true;
            this.modpacks.push(modpack);
        }
    }


    angular.module("ModpackHelper").service("application", Application);


    electron.ipcRenderer.on("update-error", (event:Electron.IpcRendererEvent, message:any) => {
        console.error(message);
    })
}
