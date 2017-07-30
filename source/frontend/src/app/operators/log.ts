import {Observable} from "rxjs/Observable";

/**
 * Logs the value of the stream the console.log everytime it emits
 *
 * @returns {Observable<T>}
 */
function log<T>(this: Observable<T>): Observable<T> {
  return this.map(value => {
    console.log(value);
    return value;
  });
}

Observable.prototype.log = log;

declare module 'rxjs/observable' {
  interface Observable<T> {
    log: typeof log;
  }
}