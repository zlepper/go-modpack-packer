import {
    app,
    BrowserWindow
} from 'electron'

import fs = require("fs");
import path = require("path");
import childprocess = require("child_process");
import {IpcHandlersCreator} from './IpcHandlers';
// There is no typings available for electron-connect, and i'm lazy, so any will have to do
//var client:any = require("electron-connect");

let win : Electron.BrowserWindow;

function startGoServer() {
    var platform = process.platform;
    let executeable: string;
    if(platform === "win32") {
        executeable = "backend.exe";
    } else {
        executeable = "backend";
    }

    // Create the backend service, and tell it where to save data
    var backend = childprocess.spawn(executeable, [app.getPath("userData")]);
    backend.stdout.on("data", function(data: any) {
        console.log(data.toString());
    });
    backend.stderr.on("data", function(data: any) {
        console.log(data.toString());
    });
}

function createWindow() {

    // Make sure to open the window where the user closed it, and with the same size
    var initPath = path.join(app.getPath("userData"), "init.json");
    var data:{bounds: Electron.Rectangle};
    try  {
        data = JSON.parse(fs.readFileSync(initPath, "utf8"));
    } catch(e) {}

    // Create the browser window
    win = new BrowserWindow((data && data.bounds) ? data.bounds : {width: 800, height: 600, frame: true});

    // and load the index.body of the app.
    win.loadURL(`file://${__dirname}/index.html`);
    win.webContents.openDevTools();
    // live reload from electron connect
    //client.client.create(win);
    
    // Save the window state, so it opens in that place next time
    win.on("close", function() {
        var data = {
            bounds: win.getBounds()
        };
        fs.writeFileSync(initPath, JSON.stringify(data));
    });

    // Emitted when the window is closed
    win.on("closed", () => {
        win = null;
    });

    IpcHandlersCreator.bindListeners();
    startGoServer()
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', createWindow);

// Quit when all windows are closed.
app.on('window-all-closed', () => {
    // On OS X it is common for applications and their menu bar
    // to stay active until the user quits explicitly with Cmd + Q
    if (process.platform !== 'darwin') {
        app.quit();
    }
});

app.on('activate', () => {
    // On OS X it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (win === null) {
        createWindow();
    }
});
