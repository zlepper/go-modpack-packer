import {ForgeVersion} from "app/models/forgeversion";
export class AWSConfig {
  public accessKey: string = "";
  public secretKey: string = "";
  public region: string = "us-east-1";
  public bucket: string = "";
}

export class FtpConfig {
  public url: string = "";
  public username: string = "";
  public password: string = "";
  public path: string = "";
}


export class UploadConfig {
  public type: string = "none";
  public aws: AWSConfig = new AWSConfig();
  public ftp: FtpConfig = new FtpConfig();
}

export class TechnicConfig {
  public isSolderPack: boolean = true;

  public createForgeZip: boolean = false;
  public forgeVersion: ForgeVersion;

  public checkPermissions: boolean = false;
  public isPublicPack: boolean = true;

  public memory: number = 0;
  public java: string = "1.8";

  public upload: UploadConfig = new UploadConfig();

  public repackAllMods: boolean = false;
}

export class SolderInfo {
  public use: boolean = false;
  public url: string = "";
  public username: string = "";
  public password: string = "";
}

export class FtbConfig {
  public isPublicPack: boolean = true;
}

export class Folder {
  public name: string;
  public include: boolean;
}

export class Modpack {
  public name: string;
  public inputDirectory: string = "";
  public outputDirectory: string = "";
  public clearOutputDirectory: boolean = true;
  public minecraftVersion: string = "1.9";
  public version: string = "1.0.0";
  public additionalFolders: Array<Folder> = [];
  public technic: TechnicConfig = new TechnicConfig();
  public ftb: FtbConfig = new FtbConfig();
  public solder: SolderInfo = new SolderInfo();
  public isNew: boolean = false;
  public id: number;

  constructor(name: string) {
    this.name = name;
    this.id = Math.floor(Math.random() * 1000000);
  }

  public static fromJson(data: Modpack): Modpack {
    let modpack = new Modpack(data.name);
    modpack.id = Math.floor(Math.random() * 1000000);
    for(let prop in data) {
      modpack[prop] = data[prop];
    }
    return modpack;
  }

  /**
   * Checks if the basic modpack info is valid.
   *
   * Returns an empty string if the basic modpack info is valid, otherwise an error message key.
   * @returns {string}
   */
  public isValid(): string {
    if (!this.name) return "Modpack name is missing.";
    if (!this.inputDirectory) return "Modpack input directory is missing.";
    if (!this.outputDirectory) return "Modpack output directory is missing.";
    if (!this.minecraftVersion) return "Modpack minecraft version has not been specified.";
    if (!this.version) return "Modpack version is missing.";
    return "";
  }

  public isValidTechnic(): string {
    var t = this.technic;
    console.log(t.upload.type);
    switch (t.upload.type) {
      case "s3":
        var aws = t.upload.aws;
        if (!aws.accessKey) {
          return "You have specified aws upload, but haven't provided an accesskey.";
        }
        if (!aws.bucket) {
          return "You have specified aws upload, but haven't specified which bucket to upload to. To select a bucket test the aws connection, and then select a bucket from the dropdown.";
        }
        if (!aws.region) {
          return "You have specified aws upload, but haven't specified a region to upload to.";
        }
        if (!aws.secretKey) {
          return "You have specified aws upload, but haven't provided a secret key.";
        }
        break;
      case "ftp":
        var ftp = t.upload.ftp;
        // Don't validate password existing, because it's possible to connect to ftp without a password
        // even though that is a very bad idea. Security and all that.
        if (!ftp.url) {
          return "You have specified ftp upload, but haven't specified a url to upload to.";
        }
        if (!ftp.username) {
          return "You have specified ftp upload, but haven't specified a username to use during connection.";
        }
        break;
      case "none":
      default:
        break;
    }
    if (t.isSolderPack) {
      var solder = this.solder;
      if (solder.use) {
        if (!solder.url) {
          return "You have not provided a url for your solder.";
        }
        if (/.*\/api\/?$/.test(solder.url)) {
          return "The solder url provided should not end with /api/, but should be the url you hit when you browser through your browser."
        }
        if (!solder.password) {
          return "You have not provided a password for your solder.";
        }
        if (!solder.username) {
          return "You have not provided a username for your solder.";
        }
      }
    }
    return "";
  }
}

export class UserPermission {
  public licenseLink: string;
  public modLink: string;
  public permissionLink: string;
  public policy: string;
  public modId: string;
}
