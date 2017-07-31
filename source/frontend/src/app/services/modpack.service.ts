import {Injectable} from "@angular/core";
import {Router} from "@angular/router";
import {Modpack} from "app/models/modpack";
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {BehaviorSubject, Observable, Subject} from "rxjs";

@Injectable()
export class ModpackService {

  private _modpacks: Subject<Modpack[]>;
  private _selectedModpack: Subject<Modpack | null>;

  constructor(protected backendCommunication: BackendCommunicationService, protected router: Router) {
    this._modpacks = new BehaviorSubject([]);
    this._selectedModpack = new BehaviorSubject(null);

    console.log('requesting load');
    backendCommunication.send('load-modpacks', {});

    backendCommunication.getMessages<Modpack[]>('data-loaded')
      .map(modpacks => modpacks.map(modpack => Modpack.fromJson(modpack)))
      .subscribe(modpacks => {
        console.log('data-loaded', modpacks);
        this._modpacks.next(modpacks);
        if (modpacks.length) {
          this._selectedModpack.next(modpacks[0]);
          this.router.navigate(['modpack']);
        }
      });
  }

  public get modpacks(): Observable<Modpack[]> {
    return this._modpacks;
  }

  public get selectedModpack(): Observable<Modpack> {
    return this._selectedModpack;
  }

  /**
   * Adds a new modpack to the modpack list
   * @param name
   */
  public addModpack(name: string) {
    const pack = new Modpack(name);
    this._modpacks.take(1).subscribe(modpacks => {
      this._modpacks.next([...modpacks, pack]);
    });
    return pack;
  }

  /**
   * Remove a given modpack from the modpacks
   * @param id
   */
  public removeModpack(id: number) {
    this._modpacks.take(1).subscribe(modpacks => {
      this._modpacks.next(modpacks.filter(modpack => modpack.id !== id));
    });
  }

  /**
   * Changes what modpack is selected
   * @param id
   */
  public setSelectedModpack(id: number) {
    this._modpacks.take(1)
      .flatMap(modpacks => modpacks)
      .filter(modpack => modpack.id === id)
      .subscribe(modpack => {
        this._selectedModpack.next(modpack)
      });
  }
}
