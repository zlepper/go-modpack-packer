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


@NgModule({
  declarations: [
    AppComponent,
    NoModpackSelectedComponent,
    SettingsComponent,
    FtbComponent,
    ModpackComponent,
    TechnicComponent,
    HeaderComponent,
    BodyComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    MaterialModule,
    FlexLayoutModule,
    routes
  ],
  providers: [
    {
      provide: APP_BASE_HREF,
      useValue: '/'
    },
    ElectronService
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
