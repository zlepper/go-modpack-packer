import {Injectable} from "@angular/core";
import {ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot} from "@angular/router";
import {ModpackService} from "app/services/modpack.service";
import {Observable} from "rxjs/Observable";

@Injectable()
export class HasModpackSelectedGuard implements CanActivate {
  constructor(protected modpackService: ModpackService) {

  }

  canActivate(next: ActivatedRouteSnapshot,
              state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.modpackService.selectedModpack.map(modpack => modpack !== null);
  }
}
