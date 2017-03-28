/// <reference path="node_modules/electron-builder/out/electron-builder.d.ts" />

const gulp = require('gulp'),
    path = require("path"),
    gutil = require('gulp-util'),
    uglify = require("gulp-uglify"),
    ts = require("gulp-typescript"),
    nodesource = "node_modules/",
    merge = require("merge-stream"),
    builder = require('electron-builder'),
    Platform = builder.Platform,
    fs = require("fs"),

    jshint = require('gulp-jshint'),
    sass = require('gulp-sass'),
    concat = require('gulp-concat'),
    sourcemaps = require('gulp-sourcemaps'),

    {exec, spawn} = require("child_process"),

    input = {
        typescript: [
            "typings/browser.d.ts",
            "source/javascript/main.ts"
        ],

        mainTypescript: [
            "typings/main.d.ts",
            "source/app/**/*.ts"
        ],

        frontendDir: 'source/frontend'
    },

    publicDir = "/public",
    defaultMainOutput = "./app",
    defaultOutputLocation = defaultMainOutput + publicDir,
    releaseBaseOutput = "./release";

function watch() {

    "use strict";
    var t4 = gulp.watch(input.mainTypescript, ["build-electron"]);
    var t5 = gulp.watch('source/backend/**/*.go', ['go-compile']);
    angularBuildWatch();
    return [t4, t5];
}

/* run the watch task when gulp is called without arguments */
gulp.task('default', [ 'build-electron', 'go-compile']);
gulp.task('build-watch', ['default'], watch);

function finishExec(callback) {
    return function (err, stdout, stderr) {
        !!stdout && console.log(stdout);
        !!stderr && console.log(stderr);
        callback(err)
    }
}

function angularBuildProd(callback) {
    return exec('ng build --prod --aot', {cwd: input.frontendDir}, finishExec(callback));
}

gulp.task('frontend-build:prod', (callback) => {
    return angularBuildProd(callback);
});

function angularBuildWatch() {
    let buildProcess = process.platform === 'win32' ?
        spawn('cmd', ['/c', 'ng build --watch'], {cwd: input.frontendDir}) :
        spawn('/bin/sh', ['-c', 'ng build --watch'], {cwd: input.frontendDir});
    buildProcess.on('close', (code) => {
        console.log(`angular watch process exited with code ${code}`);
    });

    buildProcess.stdout.on('data', (data) => {
        console.log(data.toString());
    });

    buildProcess.stderr.on('data', (data) => {
        console.error(data.toString());
    });

    return buildProcess;
}

gulp.task('frontend-build:watch', () => {
    return angularBuildWatch();
});

function electron(output) {
    if (!output) {
        output = defaultMainOutput;
    }

    const tsElectronProject = ts.createProject("source/app/tsconfig.json");
    var tsTask = gulp.src(input.mainTypescript)
        .pipe(sourcemaps.init())
        .pipe(tsElectronProject())
        .pipe(sourcemaps.write())
        .pipe(gulp.dest(output));
    var packageTask = gulp.src("source/app/package.json")
        .pipe(gulp.dest(output));
    return merge(tsTask, packageTask);
}
gulp.task("electron", function () {
    "use strict";
    return electron();
});

function buildElectron(cb, output) {
    if (!output) {
        output = defaultMainOutput;
    }
    exec("npm install", {cwd: output}, finishExec(cb));
}
gulp.task('build-electron', ["electron"], function (cb) {
    return buildElectron(cb);
});

gulp.task("go-test", function (cb) {
    exec("go test ./source/backend/...", finishExec(cb))
});

function goCompile(os, arch, cb, output) {
    if (!output) {
        output = defaultMainOutput;
    }
    // Calculate build command
    var command = "go build -o ";
    var outputFileName;
    if (os === "windows") {
        outputFileName = "backend.exe";
    } else {
        outputFileName = "backend";
    }
    command += path.join(output, outputFileName) + " ./source/backend";

    // Create special environment for cross compilation
    var env = JSON.parse(JSON.stringify(process.env));
    env.GOOS = os;
    env.GOARCH = arch;

    exec(command, {env: env}, finishExec(cb));
}

gulp.task("go-compile", function (cb) {
    var arch = process.arch === 'x64' ? "amd64" : "386";
    if (process.platform === "win32") {
        return goCompile("windows", arch, cb);
    } else if (process.platform === 'darwin') {
        return goCompile("darwin", arch, cb);
    } else {
        return goCompile("linux", arch, cb)
    }
});

function buildRelease(os, arch, callback) {
    var folder = path.join(releaseBaseOutput, os + "-" + arch);
    var tasks = [];

    // Application typescript
    tasks.push(new Promise(function (resolve) {
        electron(folder).on('finish', resolve);
    }));

    // Electron dependencies
    tasks.push(new Promise(function (resolve) {
        buildElectron(resolve, folder);
    }));

    // Go build
    tasks.push(new Promise(function (resolve) {
        goCompile(os, arch, resolve, folder);
    }));

    tasks.push(new Promise(function(resolve) {
        angularBuildProd(resolve);
    }));

    Promise.all(tasks).then(function() {
        var target;
        switch (os) {
            case 'darwin':
                target = Platform.MAC.createTarget(null, builder.Arch.x64);
                break;
            case 'windows':
                target = Platform.WINDOWS.createTarget(null, arch === '386' ? builder.Arch.ia32 : builder.Arch.x64);
                break;
            case 'linux':
                target = Platform.LINUX.createTarget(null, arch === '386' ? builder.Arch.ia32 : builder.Arch.x64);
        }
        return builder.build({
            targets: target,
            devMetadata: {
                build: {
                    appId: "dk.zlepper.modpackpacker",
                    "app-category-type": "public.app-category.developer-tools",
                    win: {
                        iconUrl: "https://raw.githubusercontent.com/zlepper/TechnicSolderHelper/master/TechnicSolderHelper/modpackhelper.ico",
                        icon: path.join(__dirname, "build", "icon.ico")
                    },
                    compression: "maximum"
                },
                directories: {
                    app: folder,
                    output: 'dist/' + os + "-" + arch
                }
            }
        });
    }).then(callback).catch(callback);
}

gulp.task("build-release:windows:x32", function (cb) {
    buildRelease('windows', '386', cb);
});

gulp.task("build-release:windows:x64", function (cb) {
    buildRelease('windows', 'amd64', cb);
});

gulp.task("build-release:linux:x64", function (cb) {
    buildRelease('linux', 'amd64', cb);
});

gulp.task('build-release:linux:x32', function(cb) {
    buildRelease('linux', '386', cb);
});

gulp.task('build-release:mac', function (cb) {
    buildRelease('darwin', 'amd64', cb);
});


gulp.task("build-release:all", [ "build-release:windows", "build-release:linux", "build-release:mac"], function () {

});

gulp.task("build-release:linux", ["build-release:linux:x64", 'build-release:linux:x32'], function () {

});

gulp.task('build-release:windows', ["build-release:windows:x32", "build-release:windows:x64"], function(){});

gulp.task('create-new-version', function(cb) {
    function incrementVersion(location, resolve, reject) {
        fs.readFile(location, "utf8", function(err, content) {
            if (err) {
                return reject(err);
            }
            var c = JSON.parse(content);
            var versionParts = c.version.split('.');
            var minor = parseInt(versionParts[2], 10);
            minor++;
            versionParts[2] = minor;
            c.version = versionParts.join('.');
            var newcontent;
            try {
                newcontent = JSON.stringify(c, null, "  ");
            } catch(err) {
                console.log(content);
                console.log(c);
                return reject(err)
            }
            fs.writeFile(location, newcontent, "utf8", function (err) {
                if (err) {
                    return reject(err);
                }
                gutil.log("New version: " + c.version);
                resolve();
            })
        });
    }
    var tasks = [];
    tasks.push(new Promise(function(resolve, reject) {
        incrementVersion("./package.json", resolve, reject);
    }));
    tasks.push(new Promise(function(resolve, reject) {
        incrementVersion("./source/app/package.json", resolve, reject);
    }));
    Promise.all(tasks).then(function() {
        cb();
    }).catch(cb);
});
