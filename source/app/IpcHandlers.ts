import {
    ipcMain,
    dialog,
    autoUpdater,
    app
} from 'electron';
import {join} from "path";
import {writeFile, readFile} from "fs";


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

        ipcMain.on("restart", () => {
            autoUpdater.quitAndInstall();
        });

        ipcMain.on("save-languages", (event:Electron.IpcMainEvent, languages: any) => {
            console.log("Saving languages");
            var folder = app.getPath("userData");
            var file = join(folder, "languages.json");

            writeFile(file, JSON.stringify(languages), {encoding: "utf8"}, function(err) {
                err && console.error(err);
                console.log("Saved languages.");
            });
        });

        ipcMain.on("get-languages", (event: Electron.IpcMainEvent) => {
            var folder = app.getPath("userData");
            var file = join(folder, "languages.json");

            readFile(file, "utf8", function(err, data) {
                if(err) {
                    event.sender.send("got-languages", "{}");
                    return console.error(err);
                }

                event.sender.send("got-languages", JSON.parse(data));
            });
        });
    }
}

