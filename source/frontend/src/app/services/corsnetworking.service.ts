import {Injectable} from "@angular/core";
import {NetworkService} from "app/services/network.service";
import {Observable} from "rxjs/Observable";

/**
 * A special service that gets around cors networking issues
 * by simply making the backend server do the request instead
 */
@Injectable()
export class CorsNetworkingService {
  constructor(private networkService: NetworkService) {

  }

  /**
   * Fetches info about a specific object from the remote
   * @param route The route to fetch from
   * @param headers The optional headers to attach to the call
   * @returns {Observable<R>}
   */
  public get<R>(route: string, headers?: { [key: string]: string }): Observable<R> {
    const backendRoute = encodeURIComponent(route);

    return this.networkService.get(`http://localhost:8084/corsaround?url=${backendRoute}`, headers);
  }
}
