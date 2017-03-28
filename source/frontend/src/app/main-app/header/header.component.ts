import {Component, OnInit} from "@angular/core";
import {Observable} from "rxjs";
import {ElectronService} from "app/services/electron.service";

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {

  constructor(protected electron: ElectronService) {
  }

  ngOnInit() {
    this.fullscreenIcon = this.electron.isMaximized
      .map(isMaximised => isMaximised ? 'fullscreen_exit' : 'fullscreen');
  }

  fullscreenIcon: Observable<string>;

  toggleMaximized() {
    this.electron.toggleMaximized();
  }

  minimize() {
    this.electron.minimize();
  }

  close() {
    this.electron.close();
  }

}
