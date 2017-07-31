import {Observable} from "rxjs/Observable";
import {Subscriber} from "rxjs/Subscriber";

/**
 *
 * @param {Observable<boolean>} pauser - An observable that pushes true when values can
 * @returns {Observable<T>}
 */
function pauseableBuffered<T>(this: Observable<T>, pauser: Observable<boolean>): Observable<T> {
  return Observable.create((subscriber: Subscriber<T>) => {
    let shouldBuffer = true;
    let updateSubscription = pauser.subscribe(s => {
      shouldBuffer = s;
      console.log('should buffer', shouldBuffer);
    });
    let buffer: Array<T> = [];
    let emptyBacklogSubscription = pauser.filter(s => !s).subscribe(s => {
      while (buffer.length) {
        let item = buffer.shift();
        subscriber.next(item);
      }
    });

    return this.subscribe((value: T) => {
        if (shouldBuffer) {
          buffer.push(value);
        } else {
          subscriber.next(value);
        }
      },
      err => subscriber.error(err),
      () => {
        emptyBacklogSubscription.unsubscribe();
        updateSubscription.unsubscribe();
        return subscriber.complete();
      }
    );
  });
}

Observable.prototype.pauseableBuffered = pauseableBuffered;

declare module 'rxjs/observable' {
  interface Observable<T> {
    pauseableBuffered: typeof pauseableBuffered;
  }
}
