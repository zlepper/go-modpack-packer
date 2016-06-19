import {
    platform,
    arch
} from "os"

import {
    app,
    BrowserWindow,
    autoUpdater
} from 'electron'

import {createReadStream, createWriteStream, readFileSync, writeFileSync} from 'fs'

import {join} from 'path';
import {spawn} from 'child_process';
import {IpcHandlersCreator} from './IpcHandlers';

(function () {
    if (platform() == "win32") {
        if (require('electron-squirrel-startup')) process.exit(0);
    }

    var devMode:boolean = (process.argv || []).indexOf('--dev') !== -1;
    let win:Electron.BrowserWindow = null;


    const shouldQuit:boolean = app.makeSingleInstance(() => {
        if (win.isMinimized()) {
            win.restore();
        }
        win.focus();
    });

    if (shouldQuit) {
        app.quit();
        return;
    }

    function canAutoupdate(): boolean {
        return !devMode && platform() === "win32";
    } 
    
    
    function setupAutoUpdater():void {
        // Don't even attempt to write the auto update unless we are on a system that supports it. 
        // Which as of the time of this comment only is windows. 
        if (!canAutoupdate()) return;
        autoUpdater.addListener("update-available", function () {
            win.webContents.send("update-info", "UPDATE.AVAILABLE");
        });
        autoUpdater.addListener("update-downloaded", function () {
            win.webContents.send("update-info", "UPDATE.DOWNLOADED");
        });
        autoUpdater.addListener("error", function (error:any) {
            console.log(error);
            win.webContents.send("update-error", error);
        });
        autoUpdater.addListener("checking-for-update", function () {
            win.webContents.send("update-info", "UPDATE.CHECKING_FOR_UPDATE");
        });
        autoUpdater.addListener("update-not-available", function () {
            win.webContents.send("update-info", "UPDATE.NOT_AVAILABLE");
        });

        var feedUrl:string = "";
        if (platform() == "win32") {
            feedUrl = "http://zlepper.dk:3215/update/win";
            if (arch() == "x64") {
                console.log("64x windows detected.");
                feedUrl += "64"
            } else {
                feedUrl += "32"
            }
        }

        feedUrl += "/" + app.getVersion();

        // The feed url was not set for some reason, so we'll not attempt to get any update package.
        if (feedUrl == "") {
            return;
        }
        autoUpdater.setFeedURL(feedUrl);
        autoUpdater.checkForUpdates();
    }

    function unpackBackend(filename:string, cb:any) {
        if (devMode) return cb();
        var read = createReadStream(join("resources", "app.asar", filename));
        var write = createWriteStream(filename);
        read.on("close", function () {
            cb();
        });
        read.pipe(write);
    }

    function startGoServer() {
        let executeable:string;
        if (platform() === "win32") {
            executeable = "backend.exe";
        } else {
            executeable = "backend";
        }

        unpackBackend(executeable, function () {
            // Create the backend service, and tell it where to save data
            var backend = spawn(executeable, [app.getPath("userData")]);
            backend.stdout.on("data", function (data:any) {
                console.log(data.toString());
            });
            backend.stderr.on("data", function (data:any) {
                console.log(data.toString());
            });
        });
    }

    function createWindow() {

        // Make sure to open the window where the user closed it, and with the same size
        var initPath = join(app.getPath("userData"), "init.json");
        var data:{bounds:Electron.Rectangle};
        try {
            data = JSON.parse(readFileSync(initPath, "utf8"));
        } catch (e) {
        }

        var bounds: Electron.BrowserWindowOptions = data && data.bounds ? data.bounds : {width: 800, height: 600};
        bounds.frame = false;

        // Create the browser window
        win = new BrowserWindow(bounds);

        // and load the index.body of the app.
        win.loadURL(`file://${__dirname}/index.html`);
        if (devMode) win.webContents.openDevTools();
        // live reload from electron connect
        //client.client.create(win);

        // Save the window state, so it opens in that place next time
        win.on("close", function () {
            var data = {
                bounds: win.getBounds()
            };
            writeFileSync(initPath, JSON.stringify(data));
        });

        // Emitted when the window is closed
        win.on("closed", () => {
            win = null;
        });

        win.webContents.on("did-finish-load", () => {
            setupAutoUpdater();
        });

        IpcHandlersCreator.bindListeners();
        startGoServer()
    }

    // This method will be called when Electron has finished
    // initialization and is ready to create browser windows.
    // Some APIs can only be used after this event occurs.
    app.on('ready', function () {
        try {
            createWindow()
        } catch (e) {
            alert(e);
        }
    });

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
})();
