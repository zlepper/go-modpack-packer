# Modpack Packer
Modpack helper for creating Technic Solder packs.
Based on electron and go.

## Build it yourself

### 1. Requirements
* [NodeJS](https://nodejs.org/)
* [Go](https://golang.org/)
* [Typings](https://github.com/typings/typings)
* [Electron-prebuilt](http://electron.atom.io/)
* [Gulp](http://gulpjs.com/)

The last 3 can be installed by running `npm install -g typings electron-prebuilt gulp`. 

### 2. Setup go
In case go is not installed _( check with `go version`)_ read up https://golang.org/doc/install
If you do not already have a setup fully working Go environment and Go workspace, read this: https://golang.org/doc/code.html#Workspace

In short, you need to setup the GOPATH variable, to some folder on your system and `GOPATH/bin` must be included in PATH.

### 3. Get the Project
**This project needs no `git clone`!**

Once you have that setup, get the repository by running `go get github.com/zlepper/go-modpack-packer/source/backend`. 
This will fetch the repo and install all go dependencies. You can now find the code in `GOPATH/src/github.com/zlepper/go-modpack-packer`

### 4. Install Dependencies
Fetch all dependencies by running `npm i`
 

### 5. Build and run
To build everything run the commmand `npm start` in the root directory of the repository. 

This will build the application and place it in the `app` directory by using **gulp** and then starts the application by running `electron . --dev` in the **app** dir. 
The `--dev` flag tells the app to run in dev mode. This way it doesn't check for updates, and also doesn't do some unpacking required when distributing. 
