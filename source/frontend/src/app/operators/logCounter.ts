import {Observable} from "rxjs/Observable";

let nextId = 0;

/**
 * Counts the number of emits going through
 *
 * @returns {Observable<T>}
 */
function logCounter<T>(this: Observable<T>): Observable<T> {
  let counts = 0;
  const id = nextId++;
  return this.do(() => {
    console.log(`id: ${id} counts: ${counts++}`);
  });
}

Observable.prototype.logCounter = logCounter;

declare module 'rxjs/observable' {
  interface Observable<T> {
    logCounter: typeof logCounter;
  }
}
