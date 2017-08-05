import {Component, OnInit} from "@angular/core";
import {ModpackService} from "app/services/modpack.service";
import {Observable} from "rxjs/Observable";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

  public ready: Observable<boolean>;

  constructor(private modpackService: ModpackService) {
  }

  ngOnInit(): void {
    this.ready = this.modpackService.ready;
  }

}
