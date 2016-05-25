/**
 * Created by Rasmus on 5/21/2016.
 */
module ModpackHelper{
    var app = angular.module("ModpackHelper", ["ngMaterial", "ui.router", "ngResource", "ngAnimate", "pascalprecht.translate", "ngMessages", "ngWebSocket"]);

    class ModpackHelperConfigure {
        static $inject = ["$mdThemingProvider", "$stateProvider", "$urlRouterProvider", "$locationProvider", "$sceDelegateProvider", "$translateProvider", "$translatePartialLoaderProvider"];
        constructor(protected $mdThemingProvider: angular.material.IThemingProvider,
                    protected $stateProvider: angular.ui.IStateProvider,
                    protected $urlRouterProvider: angular.ui.IUrlRouterProvider,
                    protected $locationProvider: angular.ILocationProvider,
                    protected $sceDelegateProvider: angular.ISCEDelegateProvider,
                    protected $translateProvider: angular.translate.ITranslateProvider,
                    protected $translatePartialLoaderProvider: angular.translate.ITranslatePartialLoaderProvider) {
            $mdThemingProvider.theme("default");
            console.log("Configuring");
            //$locationProvider.html5Mode(true);

            $urlRouterProvider.otherwise("home");

            $stateProvider.state("home", {
                url: "/",
                templateUrl: "parts/home.html",
                controller: "HomeController",
                controllerAs: "home"
            }).state("technic", {
                url: "/technic",
                templateUrl: "parts/technic.html",
                controller: "TechnicController",
                controllerAs: "vm"
            }).state("modpack", {
                url: "/modpack",
                templateUrl: "parts/modpack.html",
                controller: "ModpackController",
                controllerAs: "vm"
            }).state("ftb", {
                url: "/ftb",
                templateUrl: "parts/ftb.html",
                controller: "FtbController",
                controllerAs: "vm"
            });

            $sceDelegateProvider.resourceUrlWhitelist(["self"]);

            $translatePartialLoaderProvider.addPart("global");
            $translateProvider.useLoader("$translatePartialLoader", {
                urlTemplate: "i18n/{part}/{lang}.json"
            });

            $translateProvider.use("en");
        }
    }

    class ConfigureTranslate {
        static $inject = ["$rootScope", "$translate"];
        constructor(protected $rootScope: angular.IRootScopeService, protected $translate: angular.translate.ITranslateService) {
            $rootScope.$on('$translatePartialLoaderStructureChanged', function() {
                $translate.refresh();
            });
        }
    }

    app.config(ModpackHelperConfigure);
    app.run(ConfigureTranslate);
    
}
