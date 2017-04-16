import {APP_BASE_HREF} from "@angular/common";
import {NgModule} from "@angular/core";
import {FlexLayoutModule} from "@angular/flex-layout";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {HttpModule} from "@angular/http";

import {BrowserModule} from "@angular/platform-browser";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {MaterialModule} from "app/collectors/material.module";
import {HasModpackSelectedGuard} from "app/guards/has-modpack-selected.guard";
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {ElectronService} from "app/services/electron.service";
import {ForgeVersionService} from "app/services/forge-version.service";
import {ModpackService} from "app/services/modpack.service";
import {NetworkService} from "app/services/network.service";
import {WebSocketService} from "app/services/websocket.service";

import {AppComponent} from "./app.component";
import {routes} from "./app.routing";
import {FtbComponent} from "./ftb/ftb.component";
import {BodyComponent} from "./main-app/body/body.component";
import {HeaderComponent} from "./main-app/header/header.component";
import {ModpackEditorComponent} from "./modpack/modpack-editor/modpack-editor.component";
import {ModpackHeaderComponent} from "./modpack/modpack-header/modpack-header.component";
import {ModpackComponent} from "./modpack/modpack.component";
import {NoModpackSelectedComponent} from "./no-modpack-selected/no-modpack-selected.component";
import {SettingsComponent} from "./settings/settings.component";
import {TechnicBuildingComponent} from "./technic/technic-building/technic-building.component";
import {FileUploadComponent} from "./technic/technic-settings/file-upload/file-upload.component";
import {FtpOptionsComponent} from "./technic/technic-settings/file-upload/ftp-options/ftp-options.component";
import {S3OptionsComponent} from "./technic/technic-settings/file-upload/s3-options/s3-options.component";
import {ForgeVersionComponent} from "./technic/technic-settings/forge-version/forge-version.component";
import {JavaVersionComponent} from "./technic/technic-settings/java-version/java-version.component";
import {PackTypeComponent} from "./technic/technic-settings/pack-type/pack-type.component";
import {SolderOptionsComponent} from "./technic/technic-settings/solder-options/solder-options.component";
import {TechnicCheckPermissionsComponent} from "./technic/technic-settings/technic-check-permissions/technic-check-permissions.component";
import {TechnicSettingsComponent} from "./technic/technic-settings/technic-settings.component";
import {TechnicComponent} from "./technic/technic.component";


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
    ModpackEditorComponent,
    PackTypeComponent,
    FileUploadComponent,
    SolderOptionsComponent,
    TechnicSettingsComponent,
    ForgeVersionComponent,
    TechnicCheckPermissionsComponent,
    FtpOptionsComponent,
    S3OptionsComponent,
    JavaVersionComponent,
    TechnicBuildingComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    MaterialModule,
    FlexLayoutModule,
    routes,
    BrowserAnimationsModule,
    ReactiveFormsModule
  ],
  providers: [
    {
      provide: APP_BASE_HREF,
      useValue: '/'
    },
    ElectronService,
    ModpackService,
    ForgeVersionService,
    NetworkService,
    BackendCommunicationService,
    WebSocketService,
    HasModpackSelectedGuard
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}

