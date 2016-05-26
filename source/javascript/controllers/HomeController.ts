module HomeController {
    export class HomeController {
        public static $inject = ["application", "$state", "forge"];
        constructor(protected application: Application.Application, protected routerService: angular.ui.IStateService, protected forge: ForgeVersion.ForgeVersionService) {
            if(application.modpacks.length > 0) {
                routerService.go("technic");
            }
        }
    }
    
    angular.module("ModpackHelper").controller("HomeController", HomeController);
}
