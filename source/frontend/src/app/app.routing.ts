import {RouterModule, Routes} from "@angular/router";
import {NoModpackSelectedComponent} from "./no-modpack-selected/no-modpack-selected.component";
import {SettingsComponent} from "./settings/settings.component";
import {TechnicComponent} from "app/technic/technic.component";
import {FtbComponent} from "app/ftb/ftb.component";
import {ModpackComponent} from "app/modpack/modpack.component";

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
    component: TechnicComponent
  },
  {
    path: 'ftb',
    component: FtbComponent
  },
  {
    path: 'modpack',
    component: ModpackComponent
  },
];

export const routes = RouterModule.forRoot(appRoutes, {useHash: true,});

