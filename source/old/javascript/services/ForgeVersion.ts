module ForgeVersion {

    import IPromise = angular.IPromise;
    class ForgeMaven {
        public number: {[id: number]:Build };
        public webpath: string;
    }

    class Build {
        public build: number;
        public Jobver: string;
        public mcversion: string;
        public version: string;
        public downloadurl: string;
        public branch: string;
    }

    export class ForgeVersion {
        public build: number;
        public downloadUrl: string;
        public minecraftVersion: string;
    }

    export class ForgeVersionService {
        static $inject = ["$http"];
        public forgeVersions: Array<ForgeVersion> = [];
        public minecraftVersions: Array<string> = [];
        constructor(protected $http: angular.IHttpService) {
            this.getForgeVersions();
        }

        public getForgeVersions() {
            var self = this;
            var p = this.$http.get("http://files.minecraftforge.net/maven/net/minecraftforge/forge/json");
            p.then(function(data: any) {
                var mavenData: ForgeMaven = data.data;
                self.buildForgeDb(mavenData);
            })
        }

        private buildForgeDb(data: ForgeMaven) {
            var concurrentGone = 0;
            var i = 1;
            while(concurrentGone < 100) {
                if(i in data.number) {
                    var mcversion = data.number[i].mcversion;
                    var version = data.number[i].version;
                    var branch = data.number[i].branch;
                    var downloadUrl: string = null;
                    downloadUrl = String.format("{0}{1}-{2}{3}/forge-{1}-{2}{3}-", data.webpath, mcversion, version, String.isNullOrWhiteSpace(branch) ? "" : "-" + branch);
                    if (i < 183)
                        downloadUrl += "client.";
                    else
                        downloadUrl += "universal.";
                    if (i < 752)
                        downloadUrl += "zip";
                    else
                        downloadUrl += "jar";

                    var fv = new ForgeVersion();
                    fv.build = data.number[i].build;
                    fv.downloadUrl = downloadUrl;
                    fv.minecraftVersion = mcversion;
                    this.forgeVersions.push(fv);

                    if (this.minecraftVersions.indexOf(mcversion) === -1) {
                        this.minecraftVersions.push(mcversion);
                    }

                    concurrentGone = 0;
                } else {
                    concurrentGone++;
                }
                i++;
            }
        }
    }

    angular.module("ModpackHelper").service("forge", ForgeVersionService)
}
