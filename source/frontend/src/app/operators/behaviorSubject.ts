import {Observable} from "rxjs/Observable";
import {BehaviorSubject} from "rxjs/BehaviorSubject";
/**
 * Turns the observable into a behavior subject, so values are remember for later
 *
 * @param {T} value - The initial value of the underlaying subject
 * @returns {Observable<T>}
 */
function behaviorSubject<T>(this: Observable<T>, value: T): Observable<T> {
  const subject = new BehaviorSubject<T>(value);

  this.subscribe(
    value => subject.next(value),
    err => subject.error(err),
    () => subject.complete()
  );

  return subject;
}

Observable.prototype.behaviorSubject = behaviorSubject;

declare module 'rxjs/observable' {
  interface Observable<T> {
    behaviorSubject: typeof behaviorSubject;
  }
}
