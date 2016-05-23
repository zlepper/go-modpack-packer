import {
    ipcMain,
    dialog
} from 'electron';

export class IpcHandlersCreator {
    constructor() {
        
    }
    
    public static bindListeners():void {
        ipcMain.on("open-input-directory-dialog", (event:Electron.IpcMainEvent) => {
                dialog.showOpenDialog({
                    properties: ["openDirectory"]
                }, function (dirs:string[]) {
                    if (dirs) {
                        event.sender.send("selected-input-directory", dirs[0]);
                    }
                });
            }
        );

        ipcMain.on("open-output-directory-dialog", (event:Electron.IpcMainEvent) => {
                dialog.showOpenDialog({
                    properties: ["openDirectory"]
                }, function (dirs:string[]) {
                    if (dirs) {
                        event.sender.send("selected-output-directory", dirs[0]);
                    }
                });
            }
        );
    }
}

