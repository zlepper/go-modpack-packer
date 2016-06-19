module SettingsController {
    
    import Language = LanguageService.Language;
    export class SettingsController {
        static $inject = ["languageService"];
        constructor(protected language: LanguageService.LanguageService) {
            
        }

        public setLanguage(language: Language):void {
            this.language.setLanguage(language);
        }
    }

    angular.module("ModpackHelper").controller("SettingsController", SettingsController);
}
