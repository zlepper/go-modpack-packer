# Modpack Packer
Modpack helper for creating Technic Solder packs.

## How to contribute

### Requirements
* [NodeJS](https://nodejs.org/)
* [Go](https://golang.org/)
* Bower
* Typings
* Electron-prebuilt
* Gulp

The last 4 can be installed by running `npm install -g typings bower electron-prebuilt gulp`. 

### Settings up Go correctly, or How to clone a Go project
If you do not already have a setup fully working Go environment and Go workspace, read this: https://golang.org/doc/code.html#Workspace

In short, you need to setup the GOPATH variable, to some folder on your system. 

Once you have the variable set, create the following directories in the GOPATH folder `src/github.com/zlepper`, do this by switching to the GOPATH directory, and running the command 

\*nix: `mkdir src/github.com/zlepper`

Windows: `mkdir src\github.com\zlepper`

Then change to the newly created directory, and clone this repository by running `git clone https://github.com/zlepper/go-modpack-packer.git`. This will create a folder called `go-modpack-packer`. 
Change into this directory

### Setup
Fetch all dependencies by running these commands:
```
npm install
bower install
typings install
```

### Build
To build everything run the gulp commmand `gulp` in the root directory of the repository. 

This will build the application and place it in the `app` directory. 

### Running
Switch to the `app` directory and run the command `electron .`. 

This will launch the application. 
