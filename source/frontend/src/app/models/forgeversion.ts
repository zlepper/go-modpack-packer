export class ForgeMaven {
  public number: {[id: number]:Build };
  public webpath: string;
}

export class Build {
  public build: number;
  public Jobver: string;
  public mcversion: string;
  public version: string;
  public downloadurl: string;
  public branch: string;
}

export class ForgeVersion {
  public build: number;
  public downloadUrl: string;
  public minecraftVersion: string;
}

