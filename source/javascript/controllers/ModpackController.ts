module ModpackController {
    function folderListContains(folders: Array<Application.Folder>, folder: string): number {
        for(var i = 0; i < folders.length; i++) {
            if(folders[i].name === folder) {
                return i;
            }
        }
        return -1;
    }

    import Folder = Application.Folder;
    export class ModpackController {
        static $inject = ["application", "electron", "$state", "$translatePartialLoader", "goComm", "$rootScope"];

        constructor(protected application:Application.Application, protected electron:Electron.ElectronService, protected $state:angular.ui.IStateService, protected $translatePartialLoader:angular.translate.ITranslatePartialLoaderService, protected goComm:GoCommService.GoCommService, protected $rootScope:angular.IRootScopeService) {
            // Get translations for this page
            var controller = this;
            $translatePartialLoader.addPart("modpack");
            $rootScope.$watch(function () {
                if(!application.modpack) {
                    return "";
                }
                return application.modpack.inputDirectory;
            }, function (newValue, oldValue) {
                if (newValue && newValue.trim())
                    goComm.send("find-additional-folders", {inputDir: newValue})
            }, true);

            // We should not show the modpack page if there isn't a selected modpack
            if (!application.modpack) {
                $state.go("home");
            }

            electron.on("selected-input-directory", (event:Electron.IpcRendererEvent, path:Array<string>) => {
                application.modpack.inputDirectory = path[0];
            });
            electron.on("selected-output-directory", (event:Electron.IpcRendererEvent, path:Array<string>) => {
                application.modpack.outputDirectory = path[0];
            });
            $rootScope.$on("found-folders", (event: angular.IAngularEvent, folders: Array<string>) => {
                // Add folders to dictionary
                var oldFolders = controller.application.modpack.additionalFolders;
                controller.application.modpack.additionalFolders = [];
                folders.forEach(function(folder) {
                    var index = folderListContains(oldFolders, folder);
                    if(index === -1) {
                        var f = new Application.Folder();
                        f.include = false;
                        f.name = folder;
                        controller.application.modpack.additionalFolders.push(f);
                    } else {
                        var f = oldFolders[index];
                        controller.application.modpack.additionalFolders.push(f);
                    }
                });

            })
        }

        public selectInputDirectory():void {
            this.electron.send('open-input-directory-dialog');
        }

        public selectOutputDirectory():void {
            this.electron.send('open-output-directory-dialog');
        }
    }

    angular.module("ModpackHelper").controller("ModpackController", ModpackController);
}
