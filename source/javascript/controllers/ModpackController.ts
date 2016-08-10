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
        static $inject = ["application", "electron", "$state", "$translatePartialLoader", "goComm", "$rootScope", "forge", "$mdDialog", "$translate"];

        constructor(protected application:Application.Application,
                    protected electron:ElectronService.ElectronService,
                    protected $state:angular.ui.IStateService,
                    protected $translatePartialLoader:angular.translate.ITranslatePartialLoaderService,
                    protected goComm:GoCommService.GoCommService,
                    protected $rootScope:angular.IRootScopeService,
                    protected forge:ForgeVersion.ForgeVersionService,
                    protected $mdDialog: angular.material.IDialogService,
                    protected $translate: angular.translate.ITranslateService) {
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

            });

            
        }

        public selectInputDirectory():void {
            this.electron.send('open-input-directory-dialog', null);
        }

        public selectOutputDirectory():void {
            this.electron.send('open-output-directory-dialog', null);
        }

        public deletePack(ev: MouseEvent):void {
            this.$translate('DETAILS.ARE_YOU_SURE_DELETE').then(t => {
                this.$translate('DETAILS.YES_DELETE').then(y => {
                    this.$translate('DETAILS.NO_DELETE').then(n => {
                        this.$mdDialog.show(
                            this.$mdDialog.confirm()
                            .textContent(t)
                            .ok(y)
                            .cancel(n)
                                .targetEvent(ev)
                        ).then(() => {
                            for(var i = 0; i < this.application.modpacks.length; i++) {
                                var mp = this.application.modpacks[i];
                                if(mp.$$hash === this.application.modpack.$$hash) {
                                    this.application.modpacks.splice(i, 1);
                                    this.application.modpack = this.application.modpacks[0];
                                    if(this.application.modpack.isNew) {
                                        this.application.modpack.isNew = false;
                                        this.application.addNewModpack();
                                    }
                                    return;
                                }
                            }
                        });
                    });
                });
            });
        }

    }

    angular.module("ModpackHelper").controller("ModpackController", ModpackController);
}
