import {Injectable} from '@angular/core';
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {Observable} from "rxjs/Observable";

export interface IFoundFolders {
  key: number;
  folders: string[];
}

@Injectable()
export class AdditionalFoldersService {

  private key: number = 0;

  constructor(private backendCommunicationService: BackendCommunicationService) {
  }

  /**
   * Searches for any additional folders to be packed in the given directory
   * @param {string} inputDir
   * @returns {Observable<string[]>}
   */
  public findAdditionalFolders(inputDir: string): Observable<string[]> {
    const key = this.key++;

    this.backendCommunicationService.send('find-additional-folders', {inputDir, key});

    return this.backendCommunicationService.getMessages<IFoundFolders>('found-folders')
      .filter(found => found.key === key)
      .take(1)
      .map(found => found.folders);
  }
}
