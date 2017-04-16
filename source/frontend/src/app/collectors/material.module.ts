import {NgModule} from "@angular/core";
import {
  MdButtonModule,
  MdCardModule,
  MdCheckboxModule,
  MdCoreModule,
  MdDialogModule,
  MdIconModule,
  MdInputModule,
  MdListModule,
  MdMenuModule,
  MdProgressBarModule,
  MdRadioModule,
  MdRippleModule,
  MdSelectModule,
  MdSidenavModule,
  MdSnackBarModule,
  MdToolbarModule
} from "@angular/material";

@NgModule({
  imports: [
    MdButtonModule,
    MdCheckboxModule,
    MdSnackBarModule,
    MdRadioModule,
    MdDialogModule,
    MdSelectModule,
    MdInputModule,
    MdSidenavModule,
    MdListModule,
    MdToolbarModule,
    MdCardModule,
    MdIconModule,
    MdMenuModule,
    MdRippleModule,
    MdProgressBarModule,
    MdCoreModule
  ],
  exports: [
    MdButtonModule,
    MdCheckboxModule,
    MdSnackBarModule,
    MdRadioModule,
    MdDialogModule,
    MdSelectModule,
    MdInputModule,
    MdSidenavModule,
    MdListModule,
    MdToolbarModule,
    MdCardModule,
    MdIconModule,
    MdMenuModule,
    MdRippleModule,
    MdProgressBarModule,
    MdCoreModule
  ]
})
export class MaterialModule {

}
