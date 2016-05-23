module TechnicController {
    export class TechnicController {
        static $inject = ["application", "$translatePartialLoader"];
        
        public forgeVersions: Array<string> = [];
        
        constructor(protected application: Application.Application, protected $translatePartialLoader: angular.translate.ITranslatePartialLoaderService) {
            $translatePartialLoader.addPart("technic");
        }
    }

    angular.module("ModpackHelper").controller("TechnicController", TechnicController);
}
