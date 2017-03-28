import {UserPermission} from "app/models/modpack";
export class Mod {
  public modid:string;
  public name:string;
  public description:string;
  public version:string;
  public mcversion:string;
  public url:string;
  public authors:string;
  public credits:string;
  public filename:string;
  public md5:string;
  // Naming is totally a hack to make sure the value does not get send to the server
  public $$isDone:boolean;
  public isOnSolder:boolean;
  public userPermission:UserPermission;
  public skip:boolean;

  public static fromJson(data:Mod):Mod {
    var m = new Mod();
    m.modid = data.modid;
    m.name = data.name;
    m.description = data.description;
    m.version = data.version;
    m.mcversion = data.mcversion;
    m.url = data.url;
    m.authors = data.authors;
    m.credits = data.credits;
    m.filename = data.filename;
    m.md5 = data.md5;
    m.isOnSolder = data.isOnSolder;
    m.userPermission = data.userPermission;
    return m;
  }

  public isValid():boolean {
    if (!this.modid) return false;
    if (!this.name) return false;
    if (!this.version) return false;
    if (!this.mcversion) return false;
    if (this.authors.length < 1) return false;

    return this.isAdvancedValid();
  }

  private isAdvancedValid():boolean {
    if (this.modid.toLowerCase().indexOf("example") > -1) {
      return false;
    }
    if (this.name.toLowerCase().indexOf("example") > -1) {
      return false;
    }
    if (this.version.toLowerCase().indexOf("example") > -1) {
      return false;
    }
    if (this.name.indexOf("${") > -1) {
      return false;
    }
    if (this.version.indexOf("${") > -1) {
      return false;
    }
    if (this.mcversion.indexOf("${") > -1) {
      return false;
    }
    if (this.modid.indexOf("${") > -1) {
      return false;
    }
    if (this.version.toLowerCase().indexOf("@version@") > -1) {
      return false;
    }
    if (this.userPermission) {
      if (this.userPermission.policy !== "Open") {
        if (!this.userPermission.licenseLink) return false;
        if (!this.userPermission.modLink) return false;
        if (!this.userPermission.permissionLink) return false;
      }
    }

    return true;
  }
}
