import {AbstractControl, NG_VALIDATORS, Validator} from "@angular/forms";
import {Directive, forwardRef} from "@angular/core";

const END_WITH_API_REGEX = /.*\/api\/?$/;

@Directive({
  selector: '[noApi][ngModel],[noApi][formControl]',
  providers: [
    { provide: NG_VALIDATORS, useExisting: forwardRef(() => NoApiValidator), multi: true }
  ]
})
export class NoApiValidator implements Validator {
  validate(c: AbstractControl): { [key: string]: any; } {
    return END_WITH_API_REGEX.test(c.value) ? {validateApi: false} : null
  }


}
