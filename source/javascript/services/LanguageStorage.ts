module LanguageStorage {
    export class LanguageStorage {
        static $inject = ["electron", "$timeout"];
        public languages:any = {};
        public hasLanguageData = false;

        constructor(protected electron:ElectronService.ElectronService, protected $timeout: angular.ITimeoutService) {
            var self = this;
            this.electron.send("get-languages", null);
            this.electron.on("got-languages", function (event:Electron.IpcRendererEvent, languages:any) {
                self.hasLanguageData = true;
                languages = languages[0];
                console.log("Got languages");
                console.log(languages);
                self.languages = languages;
            });
        }

        public getCurrentLanguage(cb:any):void {
            if (this.hasLanguageData) {
                return cb(this.get("NG_TRANSLATE_LANG_KEY"));
            }
            console.log("Waiting");
            this.$timeout(this.getCurrentLanguage.bind(this, cb), 50);
        }

        public put(name:string, value:string) {
            this.languages[name] = value;
            this.electron.send("save-languages", this.languages);
        }

        public get(name:string) {
            console.log("Fetching thing");
            return this.languages[name];
        }
    }

    angular.module("ModpackHelper").factory("languageStorage", ["electron", "$timeout", function (electron:ElectronService.ElectronService, $timeout: angular.ITimeoutService) {
        return new LanguageStorage(electron, $timeout);
    }]);
}
