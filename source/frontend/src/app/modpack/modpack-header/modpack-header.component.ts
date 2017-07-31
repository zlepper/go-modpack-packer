import {Component, EventEmitter, OnInit, Output} from '@angular/core';

@Component({
  selector: 'app-modpack-header',
  templateUrl: './modpack-header.component.html',
  styleUrls: ['./modpack-header.component.scss']
})
export class ModpackHeaderComponent implements OnInit {

  @Output()
  public save: EventEmitter<void> = new EventEmitter<void>(true);

  constructor() { }

  ngOnInit() {

  }



}
