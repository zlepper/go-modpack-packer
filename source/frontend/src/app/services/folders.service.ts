import {Injectable} from "@angular/core";
import {BackendCommunicationService} from "app/services/backend-communication.service";
import {Observable} from "rxjs/Observable";

export interface IFolderResponse {
  folders: Array<string>;
  key: number;
}

@Injectable()
export class FolderService {
  /**
   * A key used for referencing between ansvers from the backend service
   */
  private key: number;

  constructor(private backendService: BackendCommunicationService) {
    this.key = 0;
  }

  /**
   * Searches for folders
   * @param folder The folder to search in
   */
  public search(folder: string): Observable<Array<string>> {
    if (folder.trim() === '') {
      folder = '/';
    }
    const key = ++this.key;

    this.backendService.send('get-folders', {folder, key});
    return this.backendService.getMessages<IFolderResponse>('got-folders')
      .filter(folders => folders.key === key)
      .take(1)
      .map(folders => folders.folders);
  }
}
