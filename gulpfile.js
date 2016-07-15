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

    jshint = require('gulp-jshint'),
    sass = require('gulp-sass'),
    concat = require('gulp-concat'),
    sourcemaps = require('gulp-sourcemaps'),

    exec = require("child_process").exec,
    ownScripts = 'source/javascript/**/*.ts',

    input = {
        sass: [
            nodesource + "angular-material/angular-material.scss",
            nodesource + "angular-material-data-table/dist/md-data-table.css",
            'source/scss/**/*.scss'
        ],

        typescript: [
            "typings/browser.d.ts",
            "source/javascript/main.ts",
            ownScripts
        ],

        mainTypescript: [
            "typings/main.d.ts",
            "source/app/**/*.ts"
        ],

        body: [
            "source/body/**/*.html",
            "source/body/**/*.json"
        ],

        vendor: [
            nodesource + "angular/angular.js",
            nodesource + "angular-animate/angular-animate.js",
            nodesource + "angular-aria/angular-aria.js",
            nodesource + "angular-messages/angular-messages.js",
            nodesource + "angular-resource/angular-resource.js",
            nodesource + "angular-sanitize/angular-sanitize.js",
            nodesource + "angular-material/angular-material.js",
            nodesource + "angular-translate/dist/angular-translate.js",
            nodesource + "angular-translate-loader-partial/angular-translate-loader-partial.js",
            nodesource + "angular-ui-router/release/angular-ui-router.js",
            nodesource + "angular-local-storage/dist/angular-local-storage.js",
            nodesource + "angular-websocket/dist/angular-websocket.js",
            nodesource + "angular-material-data-table/dist/md-data-table.js",
            nodesource + "gsap/src/uncompressed/TweenLite.js",
            nodesource + "gsap/src/uncompressed/plugins/CSSPlugin.js"
        ]
    },

    publicDir = "/public",
    defaultMainOutput = "./app",
    defaultOutputLocation = defaultMainOutput + publicDir,
    releaseBaseOutput = "./release";
function watch() {

    "use strict";
    var t1 = gulp.watch(ownScripts, ['build-ts']);
    var t2 = gulp.watch(input.sass, ['build-css']);
    var t3 = gulp.watch(input.body, ["copy-body"]);
    var t4 = gulp.watch(input.mainTypescript, ["build-electron"]);
    var t5 = gulp.watch('source/backend/**/*.go', ['go-compile']);
    //var t6 = gulp.watch('app/**/*', [electron.restart("--enable-logging")]);
    return [t1, t2, t3, t4, t5/*, t6*/];
}

/* run the watch task when gulp is called without arguments */
gulp.task('default', ['build-css', 'vendor-js', 'build-ts', 'build-electron', 'copy-body', "go-compile"]);
gulp.task('build-watch', ['default'], watch);

/* compile scss files */
function buildCss(output) {
    if (!output) {
        output = defaultOutputLocation;
    }
    return gulp.src(input.sass)
        .pipe(sourcemaps.init())
        .pipe(sass())
        .pipe(concat("bundle.css"))
        .pipe(sourcemaps.write())
        .pipe(gulp.dest(output));
}

gulp.task('build-css', function () {
    "use strict";
    return buildCss();
});

/* concat javascript files, minify if --type production */
function buildTs(output) {
    if (!output) {
        output = defaultOutputLocation;
    }
    const tsClientProject = ts.createProject("source/javascript/tsconfig.json");
    return gulp.src(input.typescript)
        .pipe(sourcemaps.init())
        .pipe(ts(tsClientProject))
        .pipe(concat('bundle.js'))
        //only uglify if gulp is ran with '--type production'
        .pipe(gutil.env.type === 'production' ? uglify() : gutil.noop())
        .pipe(sourcemaps.write())
        .pipe(gulp.dest(output));
}
gulp.task('build-ts', function () {
    "use strict";
    return buildTs();
});

function electron(output) {
    if (!output) {
        output = defaultOutputLocation;
    }

    const tsElectronProject = ts.createProject("source/app/tsconfig.json");
    var tsTask = gulp.src(input.mainTypescript)
        .pipe(sourcemaps.init())
        .pipe(ts(tsElectronProject))
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
    exec("npm install", {cwd: output}, function (err, stdout, stderr) {
        !!stdout && console.log(stdout);
        !!stderr && console.log(stderr);
        cb(err)
    });
}
gulp.task('build-electron', ["electron"], function (cb) {
    return buildElectron(cb);
});

function copyBody(output) {
    if (!output) {
        output = defaultMainOutput;
    }
    return gulp.src(input.body).pipe(gulp.dest(output));
}
gulp.task("copy-body", function () {
    return copyBody();
});

gulp.task("go-test", function (cb) {
    exec("go test ./source/backend/...", function (err, stdout, stderr) {
        !!stdout && console.log(stdout);
        !!stderr && console.log(stderr);
        cb(err)
    })
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

    exec(command, {env: env}, function (err, stdout, stderr) {
        !!stdout && console.log(stdout);
        !!stderr && console.log(stderr);
        cb(err)
    });

}
gulp.task("go-compile", ["go-test"], function (cb) {
    // var command;
    // if (process.platform === "win32") {
    //     command = "go build -o ./app/backend.exe ./source/backend"
    // } else {
    //     command = "go build -o ./app/backend ./source/backend";
    // }
    //
    var arch = process.arch === 'x64' ? "amd64" : "386";
    if (process.platform === "win32") {
        return goCompile("windows", arch, cb);
    } else if (process.platform === 'darwin') {
        return goCompile("darwin", arch, cb);
    } else {
        return goCompile("linux", arch, cb)
    }
});

function vendorJs(output) {
    if (!output) {
        output = defaultOutputLocation;
    }
    return gulp.src(input.vendor)
        .pipe(sourcemaps.init())
        .pipe(concat('vendor.js'))
        //only uglify if gulp is ran with '--type production'
        .pipe(gutil.env.type === 'production' ? uglify() : gutil.noop())
        .pipe(sourcemaps.write())
        .pipe(gulp.dest(output));
}
gulp.task('vendor-js', function () {
    "use strict";
    return vendorJs();
});

/* Watch these files for changes and run the task on update */
gulp.task('watch', watch);

function buildRelease(os, arch, callback) {
    var folder = path.join(releaseBaseOutput, os + "-" + arch);
    var tasks = [];
    // Css build
    tasks.push(new Promise(function (resolve) {
        buildCss(folder + publicDir).on("finish", resolve);
    }));

    // Vendor js build
    tasks.push(new Promise(function (resolve) {
        vendorJs(folder + publicDir).on('finish', resolve);
    }));

    // Frontend typescript
    tasks.push(new Promise(function (resolve) {
        buildTs(folder + publicDir).on('finish', resolve);
    }));

    // Application typescript
    tasks.push(new Promise(function (resolve) {
        electron(folder).on('finish', resolve);
    }));

    // Electron dependencies
    tasks.push(new Promise(function (resolve) {
        buildElectron(resolve, folder);
    }));

    // Copy body
    tasks.push(new Promise(function (resolve) {
        copyBody(folder).on('finish', resolve);
    }));

    // Go build
    tasks.push(new Promise(function (resolve) {
        goCompile(os, arch, resolve, folder);
    }));

    Promise.all(tasks).then(function() {
        // return new Promise(function(resolve, reject) {
        //     exec("build", {cwd: folder}, function (err, stdout, stderr) {
        //         !!stdout &&  console.log(stdout);
        //         !!stderr && console.log(stderr);
        //         if(err) {
        //             return reject(err);
        //         }
        //         return resolve();
        //     });
        // });
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


gulp.task("build-release:all", [ "build-release:windows", "build-release:linux"], function () {

});

gulp.task("build-release:linux", ["build-release:linux:x64",
    'build-release:linux:x32'], function () {

});

gulp.task('build-release:windows', ["build-release:windows:x32",
    "build-release:windows:x64"], function(){});
