import { ErrorDetail } from "@/api";

/**
 * Validation errors indexed by object paths
 */
export type Errors<Paths extends string> = Partial<{
  [K in Paths]: string
}>

/**
 * Index error messages by location
 */
export function indexErrors<Paths extends string>(errors: ErrorDetail[]) {
  return Object.fromEntries(
    errors.flatMap(({ location, message }) =>
      location === undefined ? [] : [[location, message]]
    )
  ) as Errors<Paths>
}