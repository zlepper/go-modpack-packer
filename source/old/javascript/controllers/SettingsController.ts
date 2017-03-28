module SettingsController {

    var electron = require("electron");

    import Language = LanguageService.Language;
    export class SettingsController {
        static $inject = ["languageService"];
        constructor(protected language: LanguageService.LanguageService) {
            
        }

        public setLanguage(language: Language):void {
            this.language.setLanguage(language);
        }

        public showDevConsole(): void {
            electron.remote.getCurrentWebContents().openDevTools('undocked');
        }
    }

    angular.module("ModpackHelper").controller("SettingsController", SettingsController);
}
