# Modpack Packer
Modpack helper for creating Technic Solder packs.

## How to contribute

### Requirements
* [NodeJS](https://nodejs.org/)
* [Go](https://golang.org/)
* Typings
* Electron-prebuilt
* Gulp

The last 3 can be installed by running `npm install -g typings electron-prebuilt gulp`. 

### Settings up Go correctly, or How to clone a Go project
If you do not already have a setup fully working Go environment and Go workspace, read this: https://golang.org/doc/code.html#Workspace

In short, you need to setup the GOPATH variable, to some folder on your system. 

Once you have that setup, get the repository by running `go get github.com/zlepper/go-modpack-packer/source/backend`. 
This will fetch the repo and install all go dependencies. 

### Setup
Fetch all dependencies by running these commands:
```
npm install
typings install
```

If you did not install the dependencies by using the `go get` command above, change to the source/backend directory and run the `go get` command again. 

### Build
To build everything run the gulp commmand `gulp` in the root directory of the repository. 

This will build the application and place it in the `app` directory. 

### Running
Switch to the `app` directory and run the command `electron . --dev`. 
The `--dev` flag tells the app to run in dev mode. This way it doesn't check for updates, and also doesn't do some unpacking required when destributing.

This will launch the application. 
