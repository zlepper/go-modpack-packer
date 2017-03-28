import {Observable} from "rxjs/Rx";
import {Http, Response, Headers, RequestOptions} from "@angular/http";
import {Injectable, isDevMode} from "@angular/core";

@Injectable()
export class NetworkService {
  private headers: {[key: string]: string} = {};

  constructor(private http: Http) {

  }

  private addAlwaysHeaders(headers: Headers): Headers {
    // Add the headers that should be on all requests
    for (let key in this.headers) {
      if (this.headers.hasOwnProperty(key)) {
        headers.set(key, this.headers[key]);
      }
    }
    return headers;
  }


  public postJsonToNonLocal<T>(body: {[key: string]: any}, url: string, headers?: {[key: string]: string}): Observable<T> {
    // Parse payloads
    let bodyPayload = JSON.stringify(body);
    let headersPayload = new Headers(headers || {});
    headersPayload.set('Content-Type', 'application/json');

    headersPayload = this.addAlwaysHeaders(headersPayload);

    let options = new RequestOptions({headers: headersPayload});

    return this.http.post(url, bodyPayload, options)
      .map(this.extractData)
      .catch(this.handleError)
      .share();
  }

  /**
   * Adds a header that will be set on all requests
   * @param key The header key
   * @param value The header value
   */
  public addHeader(key: string, value: string): void {
    this.headers[key] = value;
  }

  /**
   * Removes a header from the always included headers.
   *
   * @param key The header to remove
   * @returns {string} The value of the removed header.
   */
  public removeHeader(key: string): string {
    var value = this.headers[key];
    delete this.headers[key];
    return value;
  }

  /**
   * Extracts the data from a network requests and makes it into a normal object
   *
   * @param res
   * @returns {any|{}}
   */
  private extractData(res: Response) {
    let body = res.json();
    return body || {};
  }

  private handleError(error: any) {
    return Observable.throw(error);
  }

  /**
   * Fetches info about a specific object from the remote
   * @param route The route to fetch from
   * @param headers The optional headers to attach to the call
   * @returns {Observable<R>}
   */
  public get<R>(route: string, headers?: {[key: string]: string}): Observable<R> {
    let headersPayload = new Headers(headers || {});

    headersPayload = this.addAlwaysHeaders(headersPayload);

    let options = new RequestOptions({headers: headersPayload});

    return this.http.get(route, options)
      .map(this.extractData)
      .catch(this.handleError)
      .share();
  }

  /**
   * Updates and object on the remote
   *
   * @param body The data of the object
   * @param route The relative url to call against
   * @param headers The optional headers to attach to the call
   */
  put(body: {[key: string]: any}, route: string, headers?: {[key: string]: string}): Observable<any> {
    // Parse payloads
    let bodyPayload = JSON.stringify(body);
    let headersPayload = new Headers(headers || {});
    headersPayload.set('Content-Type', 'application/json');

    headersPayload = this.addAlwaysHeaders(headersPayload);

    let options = new RequestOptions({headers: headersPayload});

    return this.http.put(route, bodyPayload, options)
      .map(this.extractData)
      .catch(this.handleError)
      .share();
  }

  /**
   * Updates and object on the remote
   *
   * @param body The data of the object
   * @param route The relative url to call against
   * @param headers The optional headers to attach to the call
   */
  patch<T>(body: {[key: string]: any}, route: string, headers?: {[key: string]: string}): Observable<T> {
    // Parse payloads
    let bodyPayload = JSON.stringify(body);
    let headersPayload = new Headers(headers || {});
    headersPayload.set('Content-Type', 'application/json');

    headersPayload = this.addAlwaysHeaders(headersPayload);

    let options = new RequestOptions({headers: headersPayload});

    return this.http.patch(route, bodyPayload, options)
      .map(this.extractData)
      .catch(this.handleError)
      .share();
  }

  /**
   * Deletes an object on the remote
   *
   * @param route The route to delete
   * @param headers The optional headers to attach to the call
   */
  delete(route: string, headers?: {[key: string]: string}): Observable<any> {
    let headersPayload = new Headers(headers || {});

    headersPayload = this.addAlwaysHeaders(headersPayload);

    let options = new RequestOptions({headers: headersPayload});

    return this.http.delete(route, options)
      .map(this.extractData)
      .catch(this.handleError)
      .share();
  }

  /**
   * Gets the object containing the headers that are added to all network calls
   * @returns {{}}
   */
  public getHeaders(): {[key: string]: string} {
    return this.headers;
  }
}
