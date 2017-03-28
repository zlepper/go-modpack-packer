import {APP_BASE_HREF} from "@angular/common";
import {NgModule} from "@angular/core";
import {FlexLayoutModule} from "@angular/flex-layout";
import {FormsModule} from "@angular/forms";
import {HttpModule} from "@angular/http";

import {MaterialModule} from "@angular/material";
import {BrowserModule} from "@angular/platform-browser";
import {ElectronService} from "app/services/electron.service";

import {AppComponent} from "./app.component";
import {routes} from "./app.routing";
import {FtbComponent} from "./ftb/ftb.component";
import {HeaderComponent} from "./main-app/header/header.component";
import {ModpackComponent} from "./modpack/modpack.component";
import {NoModpackSelectedComponent} from "./no-modpack-selected/no-modpack-selected.component";
import {SettingsComponent} from "./settings/settings.component";
import {TechnicComponent} from "./technic/technic.component";
import { BodyComponent } from './main-app/body/body.component';
import {ModpackService} from "app/services/modpack.service";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {ForgeVersionService} from "app/services/forge-version.service";
import {NetworkService} from "app/services/network.service";
import { ModpackHeaderComponent } from './modpack/modpack-header/modpack-header.component';
import { ModpackEditorComponent } from './modpack/modpack-editor/modpack-editor.component';


@NgModule({
  declarations: [
    AppComponent,
    NoModpackSelectedComponent,
    SettingsComponent,
    FtbComponent,
    ModpackComponent,
    TechnicComponent,
    HeaderComponent,
    BodyComponent,
    ModpackHeaderComponent,
    ModpackEditorComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    MaterialModule,
    FlexLayoutModule,
    routes,
    BrowserAnimationsModule
  ],
  providers: [
    {
      provide: APP_BASE_HREF,
      useValue: '/'
    },
    ElectronService,
    ModpackService,
    ForgeVersionService,
    NetworkService
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}

