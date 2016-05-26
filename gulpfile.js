var gulp = require('gulp'),
    path = require("path"),
    gutil = require('gulp-util'),
    uglify = require("gulp-uglify"),
    ts = require("gulp-typescript"),
    tsClientProject = ts.createProject("source/javascript/tsconfig.json"),
    tsElectronProject = ts.createProject("source/app/tsconfig.json"),
    bowersource = "bower_components/",

    jshint = require('gulp-jshint'),
    sass = require('gulp-sass'),
    concat = require('gulp-concat'),
    sourcemaps = require('gulp-sourcemaps'),

    electron = require("electron-connect").server.create({
        path: ".",
        spawnOpt: {
            cwd: path.join(__dirname, "app")
        },
        verbosy: true
    }),
    exec = require("child_process").execSync,
    ownScripts = 'source/javascript/**/*.ts',

    input = {
        sass: [
            bowersource + "angular-material/angular-material.scss",
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
            bowersource + "angular/angular.js",
            bowersource + "angular-animate/angular-animate.js",
            bowersource + "angular-aria/angular-aria.js",
            bowersource + "angular-messages/angular-messages.js",
            bowersource + "angular-resource/angular-resource.js",
            bowersource + "angular-material/angular-material.js",
            bowersource + "angular-translate/angular-translate.js",
            bowersource + "angular-translate-loader-partial/angular-translate-loader-partial.js",
            bowersource + "angular-ui-router/release/angular-ui-router.js",
            bowersource + "angular-local-storage/dist/angular-local-storage.js",
            bowersource + "angular-websocket/dist/angular-websocket.js",
            "source/javascript/vendor/TweenLite.js",
            "source/javascript/vendor/CSSPlugin.js"
        ]
    },

    output = "app/public",
    mainOutput = "app";


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
gulp.task('build-watch', ['default', "run-electron"], watch);

/* compile scss files */
gulp.task('build-css', function () {
    "use strict";
    return gulp.src(input.sass)
        .pipe(sourcemaps.init())
        .pipe(sass())
        .pipe(concat("bundle.css"))
        .pipe(sourcemaps.write())
        .pipe(gulp.dest(output));
});

/* concat javascript files, minify if --type production */
gulp.task('build-ts', function () {
    "use strict";
    return gulp.src(input.typescript)
        .pipe(sourcemaps.init())
        .pipe(ts(tsClientProject))
        .pipe(concat('bundle.js'))
        //only uglify if gulp is ran with '--type production'
        .pipe(gutil.env.type === 'production' ? uglify() : gutil.noop())
        .pipe(sourcemaps.write())
        .pipe(gulp.dest(output));
});

gulp.task('build-electron', function () {
    "use strict";
    var tsTask = gulp.src(input.mainTypescript)
        .pipe(sourcemaps.init())
        .pipe(ts(tsElectronProject))
        //only uglify if gulp is ran with '--type production'
        .pipe(gutil.env.type === 'production' ? uglify() : gutil.noop())
        .pipe(sourcemaps.write())
        .pipe(gulp.dest(mainOutput));
    var packageTask = gulp.src("source/app/package.json")
        .pipe(gulp.dest("app"));
    return [tsTask, packageTask];
});

gulp.task("copy-body", function () {
    return gulp.src(input.body).pipe(gulp.dest(mainOutput));
});

gulp.task('run-electron', function () {
    //electron.start("--enable-logging");
});

gulp.task("go-compile", function () {
    var command;
    if (process.platform === "win32") {
        command = "go build -o ./app/backend.exe ./source/backend"
    } else {
        command = "go build -o ./app/backend ./source/backend";
    }
    console.log(__dirname);
    gutil.log(command);
    return exec(command);
});

gulp.task('vendor-js', function () {
    "use strict";
    return gulp.src(input.vendor)
        .pipe(sourcemaps.init())
        .pipe(concat('vendor.js'))
        //only uglify if gulp is ran with '--type production'
        .pipe(gutil.env.type === 'production' ? uglify() : gutil.noop())
        .pipe(sourcemaps.write())
        .pipe(gulp.dest(output));
});

/* Watch these files for changes and run the task on update */
gulp.task('watch', watch);
