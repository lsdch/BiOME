import { BaseValidation, Validation, ValidationArgs, ValidationRuleCollection } from "@vuelidate/core";

// BaseValidation<T, ValidationRuleCollection<T> | undefined>
export function vuelidateErrors(item: any): string[] {
  return item.$errors.map((e: any) => String(e.$message))
}