import {RouterModule, Routes} from "@angular/router";
import {FtbComponent} from "app/ftb/ftb.component";
import {HasModpackSelectedGuard} from "app/guards/has-modpack-selected.guard";
import {ModpackComponent} from "app/modpack/modpack.component";
import {TechnicComponent} from "app/technic/technic.component";
import {NoModpackSelectedComponent} from "./no-modpack-selected/no-modpack-selected.component";
import {SettingsComponent} from "./settings/settings.component";

const appRoutes: Routes = [
  {
    path: '',
    component: NoModpackSelectedComponent
  },
  {
    path: 'settings',
    component: SettingsComponent
  },
  {
    path: 'technic',
    component: TechnicComponent,
    canActivate: [
      HasModpackSelectedGuard
    ]
  },
  {
    path: 'ftb',
    component: FtbComponent,
    canActivate: [
      HasModpackSelectedGuard
    ]
  },
  {
    path: 'modpack',
    component: ModpackComponent,
    canActivate: [
      HasModpackSelectedGuard
    ]
  },
];

export const routes = RouterModule.forRoot(appRoutes, {useHash: true,});

