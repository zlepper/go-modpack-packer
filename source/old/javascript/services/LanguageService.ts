module LanguageService {
    export class Language  {
        public key: string;
        public name: string;

        constructor(k: string, n: string) {
            this.key = k;
            this.name = n;
        }
    }


    export class LanguageService {
        static $inject = ["$translate", "$translatePartialLoader", "languageStorage"];

        public languages: Array<Language> = [
            new Language("en", "English"),
            new Language("da", "Danish")
        ];

        public language: Language = this.languages[0];

        constructor(protected $translate: angular.translate.ITranslateService,
                    protected $translatePartialLoader: angular.translate.ITranslatePartialLoaderService,
                    protected languageStorage: LanguageStorage.LanguageStorage) {
            $translatePartialLoader.addPart("settings");
            var self = this;
            languageStorage.getCurrentLanguage(function(lang: string) {
                if (lang) {
                    for (var i = 0; i < self.languages.length; i++) {
                        if (self.languages[i].key === lang) {
                            console.log("Found key");
                            self.language = self.languages[i];
                            self.setLanguage(self.language);
                            break;
                        }
                    }
                }
            });
        }

        public setLanguage(language: Language):void {
            this.$translate.use(language.key)
        }
    }

    angular.module("ModpackHelper").service("languageService", LanguageService);
}
