import { ErrorDetail, ErrorModel } from "@/api"
import * as Schemas from "@/api/gen/schemas.gen"
import { useErrorHandler, ResponseBody } from "@/api/responses"
import { IndexedValidationErrors } from "@/functions/mutations"
import { useCountries } from "@/stores/countries"
import { OpenApi } from "@hey-api/openapi-ts"
import { List, Union } from "ts-toolbelt"
import { Reactive, reactive, ref, watch } from "vue"


/**
 * Validation rules for a field, based on the OpenAPI schema.
 * Compatible with Vuetify.
 */
type Rule = ((v: any) => true | string)


type SchemaModule = typeof import("@/api/gen/schemas.gen")

type SchemaRefs = keyof SchemaModule extends `$${infer U}` ? `#/components/schemas/${U}` : never

type SchemaIndex = {
  [K in SchemaRefs]: K extends `#/components/schemas/${infer U}`
  ? `$${U}` extends keyof SchemaModule ?
  SchemaModule[`$${U}`] : never : never
}

export type Schema = NonNullable<NonNullable<OpenApi.V3_1_X['components']>['schemas']>[string]

export type SchemaProperties = Readonly<Record<string, Schema>>
export type SchemaWithProperties<P> = Schema & Readonly<{ type: "object", properties: P }>


/**
 * All property paths in an OpenAPI schema
 */

export type SchemaPaths<T extends Schema, Terminal extends "Terminal" | "All" = "All"> =
  T extends { properties: Record<string, Schema> } ? {
    [K in keyof T['properties']]-?:
    // If terminal is "All", include any property whatsoever
    | ("Terminal" extends Terminal ? never : [K])
    // Resolves to all terminal properties
    // i.e. :
    // - actual terminal property if trailing path is empty
    // - recurse if trailing path is not empty
    | [K, ...SchemaPaths<T['properties'][K], Terminal>]
  }[keyof T['properties']]
  : T extends { items: Schema }
  ? ([number, ...SchemaPaths<T['items'], Terminal>])
  : T extends { $ref: `#/components/schemas/${infer Z}` }
  ? (`$${Z}` extends keyof SchemaModule
    ? SchemaPaths<SchemaModule[`$${Z}`] extends Schema ? SchemaModule[`$${Z}`] : never, Terminal>
    : [])
  : []

/**
 * Property path derived from array path
 * @example
 * ['items', 1, 'name'] => 'items[1].name'
 * @deprecated Not actually used anywhere at the moment
 */
export type StringPath<P extends (string | number)[]> =
  // Terminal nodes
  P extends [] ? '' :
  P extends [number] ? `[${P[number]}]` :
  P extends [string] ? `${P[number]}` :
  // Non-terminal nodes
  P extends [number, ...infer Rest extends (string | number)[]] ? (
    `[${P[0]}]${Rest[0] extends string ? '.' : ''}${StringPath<Rest>}`
  ) :
  P extends [string, ...infer Rest extends (string | number)[]] ? (
    `${P[0]}${Rest[0] extends string ? '.' : ''}${StringPath<Rest>}`
  ) :
  never

/**
 * Property path definition in a schema, replacing array elements with '*'
 * @example ['items', '*', 'name']
 */
type CollectedPath<T extends Schema> = List.Replace<SchemaPaths<T, "All">, number, '*'>

/**
 * Traverse a schema definition, gathering all property paths.
 * Array elements are represented by '*' in the path.
 * @example ['items', '*', 'name']
 */
function collectPaths<T extends Schema>(s: T): Union.ListOf<CollectedPath<T>> {
  let paths = [] as Array<CollectedPath<T>>
  if (s.properties) {
    paths = Object.entries(s.properties as Record<string, Schema>)
      .reduce<Array<CollectedPath<T>>>((acc, [prop, schema]) => {
        if (schema.$ref) {
          const ref = getSchemaRef(schema.$ref as SchemaRefs)
          const p = collectPaths(ref as Schema).map(p => [prop, ...p])
          return acc.concat(p.length ? p as Array<CollectedPath<T>> : [[prop]] as [CollectedPath<T>])
        }
        if (schema.properties)
          return acc.concat(collectPaths(schema).map(p => [prop, ...p] as CollectedPath<T>))
        if (schema.items)
          return acc.concat(collectPaths(schema).map(p => [prop, '*', ...p] as unknown as CollectedPath<T>))
        return acc.concat([[prop]] as [CollectedPath<T>])
      },
        []
      )
  }
  return paths as unknown as Union.ListOf<CollectedPath<T>>
}





export type FieldSpecification = {
  schema: Schema
  required: boolean
}

/**
 * Retrieve a schema by its reference from the OpenAPI spec
 */
export function getSchemaRef<R extends SchemaRefs>(ref: R) {
  const refName = `$${ref.split('/').at(-1)}` as (keyof SchemaModule)
  return Schemas[refName] as SchemaIndex[R]
}

export function getSchema<T extends Schema>(schema: T, ...path: SchemaPaths<T, any>): FieldSpecification {

  // If schema is a reference, resolve it
  if (schema.$ref !== undefined) {
    const target = getSchemaRef(schema.$ref as SchemaRefs) as Schema
    return getSchema(target, ...path as SchemaPaths<typeof target, any>)
  }
  // Extract next path element
  const [fragment, ...rest] = path
  // Path definition ends, but matching property might not be terminal
  if (rest.length == 0) {
    // key is string: reached property definition
    if (typeof fragment === "string") {
      if (!schema.properties) {
        throw {
          error: new Error(`Expected properties in schema, attempting to access ${fragment}`),
          schema
        }
      }
      const prop = schema.properties[fragment] as Schema
      return {
        schema: prop?.$ref !== undefined ? getSchemaRef(prop.$ref as SchemaRefs) as Schema : prop,
        required: schema.required?.includes(fragment) ?? false
      }
    }
    // key is number: reached array item definition
    else if (typeof fragment === "number" && !!schema.items)
      return getSchema(schema.items)
  }
  // Non-terminal path
  if (typeof fragment === "string" && !!schema.properties)
    return getSchema(schema.properties[fragment] as Schema, ...rest)
  else if (typeof fragment === "number" && schema.items)
    return getSchema(schema.items, ...rest)
  else throw {
    error: new Error(`Invalid path fragment: ${fragment} with type ${typeof fragment}`),
    schema
  }
}

export function joinPath<T extends Schema>(path: SchemaPaths<T, "All">) {
  return path.reduce((acc: string, p) => {
    let suffix = String(p)
    if (acc.length !== 0 && typeof p === 'string') {
      suffix = `.${suffix}`
    } else if (typeof p === "number") {
      suffix = `[${suffix}]`
    }
    return `${acc}${suffix}`
  }, '')
}

export type ErrorBinding = { errorMessages?: string[] | undefined }
export type SchemaBinding = ErrorBinding & FieldBinding
export type SchemaBinder<T extends Schema> = (...path: SchemaPaths<T, "All">) => SchemaBinding

/**
 * Schema binding for form elements. Sets field constraints, client-side validation rules, hints, and classes.
 */
export type FieldBinding = {
  hint?: string,
  min?: number,
  max?: number,
  minLength?: number,
  maxLength?: number,
  rules: ((value: any) => true | string)[]
  class?: string | Record<string, boolean> | (string | Record<string, boolean>)[]
}

export function patternRule(pattern: string, errMessage = "Invalid format") {
  const regex = new RegExp(pattern)
  return (value: string) => {
    return !value || regex.test(value) || errMessage
  }
}

export type ValidationError = ErrorDetail & { path: Array<string | number> }

export type PathPrefix<T extends Schema> = Exclude<SchemaPaths<T, "All">, SchemaPaths<T, "Terminal">>

export type PathComplement<T extends Schema, Pref extends PathPrefix<T> | undefined, Paths extends SchemaPaths<T, "Terminal"> = SchemaPaths<T, "Terminal">> =
  undefined extends Pref
  ? Paths
  : Paths extends [...NonNullable<Pref>, ...infer R] ? R : never;

export type PathJoin<T extends Schema, P extends PathPrefix<T>, Complement extends PathComplement<T, P> = PathComplement<T, P>> = [...P, ...Complement] extends SchemaPaths<T, "All"> ? ([...P, ...Complement] & SchemaPaths<T, "All">) : never

function joinPathPrefix<
  T extends Schema,
  Prefix extends PathPrefix<T> = PathPrefix<T>,
>(prefix: Prefix, path: PathComplement<T, Prefix>): PathJoin<T, Prefix> {
  return [...prefix, ...path] as PathJoin<T, Prefix>
}


export function useSchema<T extends Schema>(
  schema: T,
  useErrors?: Reactive<Record<string, string[]>>
) {

  /**
   * Traverse schema, collecting all property paths.
   * Used for explicit validation calls using `validateAll`.
   */
  const allPaths = collectPaths(schema)

  /**
 * Input validation errors indexed by their object path in the API request body
 */
  const errors = ref<IndexedValidationErrors>({ "rest": [] })
  const unindexedErrors = ref<string[]>([])

  watch(
    () => useErrors,
    (e) => {
      if (e) {
        Object.entries(e).forEach(([key, value]) => {
          errors.value[key] = value
        })
      }
    })

  /**
   * Derives validation rules from an OpenAPI schema.
   * Compatible with Vuetify form validation API.
   */
  function makeRules({ schema: s, required }: FieldSpecification) {
    const rules: Rule[] = []
    if (required) rules.push((value: any) => !!value || value === 0 ? true : "This field is required")

    // Length validation
    if (s?.minLength !== undefined) {
      rules.push((value: string | Array<unknown>) => (value?.length ?? 0) >= (s.minLength!)
        || `At least ${s.minLength!} ${s.type == "string" ? 'character(s)' : 'element(s)'} required`)
    }
    if (s?.maxLength !== undefined) {
      rules.push((value: string | Array<unknown>) => (value?.length ?? 0) <= (s.maxLength!)
        || `At most ${s.maxLength!} ${s.type == "string" ? 'character(s)' : 'element(s)'} accepted`)
    }

    // Numbers
    if ((s?.type == 'number' || s?.type == "integer")) {
      rules.push((value?: string | number) => {
        if (value === undefined || value === null || value === "") return true
        value = Number(value)
        if (s.type == "integer") return Number.isInteger(value) || `Must be an integer number`
        if (s.type == "number" && s.format == "float") return (Number.isFinite(value)) || `Must be a decimal number`
        return Number.isFinite(value) || `Must be a number`
      })
    }
    if (s?.maximum !== undefined) {
      rules.push((value: number) => (value === undefined || value === null) || (value <= s.maximum!) || `Maximum value is ${s.maximum!}`)
    }
    if (s?.minimum !== undefined) {
      rules.push((value: number) => (value === undefined || value === null) || (value >= s.minimum!) || `Minimum value is ${s.minimum!}`)
    }

    // Enum
    if (s?.enum !== undefined) {
      rules.push((value: any) => !value || s.enum?.includes(value) || 'Invalid value')
    }

    // Regex
    if (s?.pattern !== undefined) {
      rules.push(patternRule(s.pattern))
    }

    // Formats
    switch (s?.format) {

      case "email":
        const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
        rules.push((value: string) => {
          return !value || emailRegex.test(value) || "Invalid email format"
        })
        break;
      // Custom
      case "country-code":
        rules.push((value: string) =>
          !value || useCountries().findCountry(value) !== undefined ||
          `Invalid country code`)
        break;

      default:
        break;
    }

    return rules
  }

  /**
   * Validates a value against the OpenAPI schema.
   * In most cases you should use `field` instead, which also binds error messages to form fields.
   * This is useful for validating values outside of form elements.
   */
  function validate(...path: SchemaPaths<T, "Terminal">) {
    const spec = getSchema(schema, ...path)
    const rules = makeRules(spec)
    return (value: any) => {
      return rules.reduce<true | string>((acc, rule) => {
        if (acc !== true) return acc
        return rule(value)
      }, true)
    }
  }



  function isArrayPath(path: CollectedPath<T> | SchemaPaths<T>): path is CollectedPath<T> {
    return (path as string[]).includes('*')
  }

  function validateAll(
    v: Array<any> | Record<string, any>,
    paths: Readonly<Array<CollectedPath<T>>> = allPaths,
    prefix?: Array<string | number>
  ): ValidationError[] {

    return paths.flatMap<ValidationError>((path: CollectedPath<T>) => {
      if (isArrayPath(path)) {
        const { errors } = path.reduce(
          (acc, p, i) => {
            if (p === '*') {
              const arrayItemPath = Array(path.slice(i + 1)) as [CollectedPath<T>]
              const arrayPathPrefix = (prefix ?? []).concat(path.slice(0, i))

              const validatedItems = (acc.value as Array<any>).flatMap((item, i) => {
                return validateAll(item, arrayItemPath, [...arrayPathPrefix, i])
              })
              return {
                errors: acc.errors.concat(validatedItems),
                value: acc.value
              }
            } else {
              return {
                value: (acc.value as Record<string, any>)[p],
                errors: acc.errors
              }
            }
          }
          ,
          {
            value: v,
            errors: new Array<ValidationError>()
          }
        )
        return errors
      }

      // Get the value at the path
      const value = path.reduce((acc, p) => (acc as Record<string, any>)[p], v)
      const valid = validate(...path)(value)
      return valid !== true ? [{
        location: path.join('.'),
        message: valid,
        value,
        path: (prefix ?? []).concat(path)
      }] : []
    })
  }

  /**
   * Generates field properties from an OpenAPI schema.
   * Compatible with Vuetify form elements.
   */
  function fieldProps(...path: SchemaPaths<T, "All">): FieldBinding {
    const spec = getSchema(schema, ...path)
    const rules = makeRules(spec)
    const { schema: s } = spec
    return {
      hint: s?.description,
      min: s?.minimum,
      max: s?.maximum,
      minLength: s?.minLength,
      maxLength: s?.maxLength,
      rules,
      class: { 'required': spec.required }
    }
  }

  /**
   * Collects error messages indexed by their object path in an API request body,
   * so that they can be consumed by `bindErrors` or `field`.
   */
  function dispatchErrors(body: ErrorModel) {
    body.errors?.forEach((e) => {
      if (e.location === undefined)
        unindexedErrors.value.push(e.message ?? "Invalid value")
      else if (e.location.startsWith('body.')) {
        const loc = e.location.replace('body.', '')
        errors.value[loc].push(e.message ?? "Invalid value")
      }
    })
  }

  /**
   * Dispatches errors from the API response, or returns response data in a promise if no errors are present.
   */
  function handleErrors<D>(e: ResponseBody<D, ErrorModel>) {
    return useErrorHandler<D, ErrorModel>(dispatchErrors)(e)
  }

  /**
   * Binds remote error messages to an input form element.
   * Errors must be caught using `errorHandler` function.
   *
   * @param path The object property path for the field
   */
  function bindErrors(...path: SchemaPaths<T, "All">): ErrorBinding {
    const strPath = joinPath(path)
    errors.value[strPath] = reactive([])
    return {
      errorMessages: errors.value[strPath]
    }
  }

  /**
   * Binds validation rules and remote error messages to an input form element,
   * using the provided OpenAPI schema.
   * Errors must be caught using `errorHandler` function.
   *
   * @example `<v-text-field v-model="model.someArray[0].someProperty" v-bind="schema('someArray', 0, 'someProperty')" />`
   * @param path The object property path for the field
   * @returns Field bindings to be passed to form element using `v-bind`
   */
  function bindSchema(...path: SchemaPaths<T, "All">): SchemaBinding {
    return {
      ...fieldProps(...path),
      ...bindErrors(...path),
    }
  }


  function withPrefix(...prefix: PathPrefix<T>) {
    const errorBinder = (...path: PathComplement<T, PathPrefix<T>>): ErrorBinding => {
      const fullPath = joinPathPrefix(prefix, path)
      return bindErrors(...fullPath)
    }
    const fieldBinder = (...path: PathComplement<T, PathPrefix<T>>): FieldBinding => {
      const fullPath = joinPathPrefix(prefix, path)
      return fieldProps(...fullPath)
    }
    return {
      error: errorBinder,
      field: fieldBinder,
      schema: (...path: PathComplement<T, PathPrefix<T>>): SchemaBinding => {
        return {
          ...errorBinder(...path),
          ...fieldBinder(...path),
        }
      }
    }
  }


  return {
    /**
     * Bind validation rules and remote error messages to an input form element,
     * using the provided OpenAPI schema.
     * Errors must be caught using `handleErrors` or `dispatchErrors` functions.
     */
    bind: {
      /**
       * Bind both client-side constraints and validation rules
       * and remote error messages to an input form element.
       */
      schema: bindSchema,
      /**
       * Unwrap nested schema and expose bindings for its properties.
       */
      withPrefix,
      /**
       * Bind client-side constraints and validation rules.
       */
      field: fieldProps,
      /**
       * Bind remote error messages to an input form element.
       * Errors must be caught using `handleErrors` or `dispatchErrors` functions.
       */
      errors: bindErrors,
    },
    /**
     * Handle errors from API responses, dispatching them to form fields.
     * Convenience wrapper around dispatchErrors that accepts a response body
     * that may or may not contain errors.
     */
    handleErrors,
    /**
     * Dispatches errors from the API response to bounds fields as error-messages prop.
     */
    dispatchErrors,
    validate,
    validateAll,
    errors,
    unindexedErrors,
    paths: collectPaths(schema),
  }
}


export function addRules<T extends { rules: Rule[] }>({ rules, ...rest }: T, ...add: Rule[]) {
  return {
    ...rest,
    rules: rules.concat(add)
  }
}