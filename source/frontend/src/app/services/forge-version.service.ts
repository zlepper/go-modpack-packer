import { Injectable } from '@angular/core';
import {ForgeMaven, ForgeVersion} from "app/models/forgeversion";
import {NetworkService} from "app/services/network.service";
import {BehaviorSubject, Observable, Subject} from "rxjs";

@Injectable()
export class ForgeVersionService {

  private _forgeVersions: Subject<ForgeVersion[]>;
  private _minecraftVersions: Subject<string[]>;

  public get forgeVersions(): Observable<ForgeVersion[]> {
    return this._forgeVersions;
  }

  public get minecraftVersions(): Observable<string[]> {
    return this._minecraftVersions;
  }

  constructor(protected networkService: NetworkService) {
    this.getForgeVersions();
    this._forgeVersions = new BehaviorSubject<ForgeVersion[]>([]);
    this._minecraftVersions = new BehaviorSubject<string[]>([]);
  }

  private isNullOrWhiteSpace(s) {
    return s === null || s.match(/^ *$/) !== null;
  }

  public getForgeVersions() {
    this.networkService.get<ForgeMaven>('http://files.minecraftforge.net/maven/net/minecraftforge/forge/json')
      .subscribe(data => {
        this.buildForgeDb(data);
      })
  }

  private buildForgeDb(data: ForgeMaven) {
    let minecraftVersions: string[] = [];
    let forgeVersions: ForgeVersion[] = [];

    let concurrentGone = 0;
    let i = 1;
    while(concurrentGone < 100) {
      if(i in data.number) {
        const mcversion = data.number[i].mcversion;
        const version = data.number[i].version;
        let branch = data.number[i].branch;
        let downloadUrl: string = null;
        branch = this.isNullOrWhiteSpace(branch) ? "" : "-" + branch;
        downloadUrl = `${data.webpath}${mcversion}-${version}${branch}/forge-${mcversion}-${version}${branch}`;
        if (i < 183)
          downloadUrl += "client.";
        else
          downloadUrl += "universal.";
        if (i < 752)
          downloadUrl += "zip";
        else
          downloadUrl += "jar";

        const fv = new ForgeVersion();
        fv.build = data.number[i].build;
        fv.downloadUrl = downloadUrl;
        fv.minecraftVersion = mcversion;
        forgeVersions.push(fv);

        if (minecraftVersions.indexOf(mcversion) === -1) {
          minecraftVersions.push(mcversion);
        }

        concurrentGone = 0;
      } else {
        concurrentGone++;
      }
      i++;
    }

    this._forgeVersions.next(forgeVersions);
    this._minecraftVersions.next(minecraftVersions);
  }
}
