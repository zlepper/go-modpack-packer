import 'rxjs/add/operator/bufferTime';
import {Observable} from "rxjs/Observable";

/**
 * Like bufferTime, but only emits when there are actually new things in the array....
 * See this issue: https://github.com/ReactiveX/rxjs/issues/2601
 *
 * @returns {Observable<T>}
 */
function betterBufferTime<T>(this: Observable<T>, time: number): Observable<T[]> {
  return this.bufferTime(time)
    .filter(array => array.length > 0);
}

Observable.prototype.betterBufferTime = betterBufferTime;

declare module 'rxjs/observable' {
  interface Observable<T> {
    betterBufferTime: typeof betterBufferTime;
  }
}
