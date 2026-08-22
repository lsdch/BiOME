import { AppError, ErrorDetail } from '@/api'
import * as Schemas from '@/api/gen/schemas.gen'
import { IndexedValidationErrors } from '@/lib/mutations'
import { useCountries } from '@/stores/countries'
import { OpenApiSchemaObject } from '@hey-api/openapi-ts'
import { List, Union } from 'ts-toolbelt'
import { computed, Ref } from 'vue'

/**
 * Validation rules for a field, based on the OpenAPI schema.
 * Compatible with Vuetify.
 */
type Rule = (v: any) => true | string

type SchemaModule = typeof import('@/api/gen/schemas.gen')

type SchemaRefs = keyof SchemaModule extends `$${infer U}` ? `#/components/schemas/${U}` : never

type SchemaIndex = {
  [K in SchemaRefs]: K extends `#/components/schemas/${infer U}`
    ? `$${U}` extends keyof SchemaModule
      ? SchemaModule[`$${U}`]
      : never
    : never
}

// export type Schema = OpenApiSchemaObject.V3_1_X
type WithReadonlyRequired<T> = T extends unknown
  ? Omit<T, 'required'> & {
      required?: readonly string[]
    }
  : never

export type Schema = WithReadonlyRequired<OpenApiSchemaObject.V3_1_X>
// export type Schema = Override<OpenApiSchemaObject.V3_1_X, { required?: readonly string[] }>
export type SchemaProperties = Readonly<Record<string, Schema>>
export type SchemaWithProperties<P> = Schema & Readonly<{ type: 'object'; properties: P }>

/**
 * All property paths in an OpenAPI schema
 */
export type SchemaPaths<T extends Schema, Terminal extends 'Terminal' | 'All' = 'All'> = T extends {
  properties: Record<string, Schema>
}
  ? {
      [K in keyof T['properties']]-?:
        // If terminal is "All", include any property whatsoever
        | ('Terminal' extends Terminal ? never : [K])
        // Resolves to all terminal properties
        // i.e. :
        // - actual terminal property if trailing path is empty
        // - recurse if trailing path is not empty
        | [K, ...SchemaPaths<T['properties'][K], Terminal>]
    }[keyof T['properties']]
  : T extends { items: Schema }
    ? [number, ...SchemaPaths<T['items'], Terminal>]
    : T extends { $ref: `#/components/schemas/${infer Z}` }
      ? `$${Z}` extends keyof SchemaModule
        ? SchemaPaths<
            SchemaModule[`$${Z}`] extends Schema ? SchemaModule[`$${Z}`] : never,
            Terminal
          >
        : []
      : []

/**
 * Property path definition in a schema, replacing array elements with '*'
 * @example ['items', '*', 'name']
 */
type CollectedPath<T extends Schema> = List.Replace<SchemaPaths<T, 'All'>, number, '*'>

function collectedPath<T extends Schema>(path: Array<string | '*'>): CollectedPath<T> {
  return path as unknown as CollectedPath<T>
}

/**
 * Traverse a schema definition, gathering all property paths.
 * Array elements are represented by '*' in the path.
 * @example ['items', '*', 'name']
 */
function collectPaths<T extends Schema>(s: T): Union.ListOf<CollectedPath<T>> {
  let paths = [] as Array<CollectedPath<T>>

  if (s.properties) {
    paths = Object.entries(s.properties as Record<string, Schema>).reduce<Array<CollectedPath<T>>>(
      (acc, [prop, schema]) => {
        if (schema.$ref) {
          const ref = getSchemaRef(schema.$ref as SchemaRefs)
          const p = collectPaths(ref as Schema).map((p) => collectedPath<T>([prop, ...p]))

          return acc.concat(p.length ? p : [collectedPath<T>([prop])])
        }

        if (schema.properties) {
          return acc.concat(collectPaths(schema).map((p) => collectedPath<T>([prop, ...p])))
        }

        if (schema.items) {
          return acc.concat(collectPaths(schema).map((p) => collectedPath<T>([prop, '*', ...p])))
        }

        return acc.concat([collectedPath<T>([prop])])
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
  const refName = `$${ref.split('/').at(-1)}` as keyof SchemaModule
  return Schemas[refName] as SchemaIndex[R]
}

function getSchemaRuntime(schema: Schema, path: readonly (string | number)[]): FieldSpecification {
  // Empty path = current schema
  if (path.length === 0) {
    return {
      schema,
      required: false
    }
  }

  // Resolve references before traversing
  if (schema.$ref !== undefined) {
    const target = getSchemaRef(schema.$ref as SchemaRefs) as Schema
    return getSchemaRuntime(target, path)
  }

  const [fragment, ...rest] = path

  // String fragment = object property
  if (typeof fragment === 'string') {
    if (!schema.properties) {
      throw {
        error: new Error(`Expected properties in schema, attempting to access ${fragment}`),
        schema
      }
    }

    const prop = schema.properties[fragment] as Schema | undefined

    if (!prop) {
      throw {
        error: new Error(`Property ${fragment} does not exist in schema`),
        schema
      }
    }

    // We reached the requested property
    if (rest.length === 0) {
      return {
        schema: prop.$ref !== undefined ? (getSchemaRef(prop.$ref as SchemaRefs) as Schema) : prop,
        required: schema.required?.includes(fragment) ?? false
      }
    }

    // Continue traversing
    return getSchemaRuntime(prop, rest)
  }

  // Number fragment = array item
  if (typeof fragment === 'number') {
    if (!schema.items) {
      throw {
        error: new Error(`Expected items in schema, attempting to access array index ${fragment}`),
        schema
      }
    }

    // We reached the array item
    if (rest.length === 0) {
      return {
        schema:
          schema.items.$ref !== undefined
            ? (getSchemaRef(schema.items.$ref as SchemaRefs) as Schema)
            : schema.items,
        required: false
      }
    }

    // Continue traversing
    return getSchemaRuntime(schema.items, rest)
  }

  throw {
    error: new Error(`Invalid path fragment: ${String(fragment)} with type ${typeof fragment}`),
    schema
  }
}

export function getSchema<T extends Schema>(
  schema: T,
  ...path: SchemaPaths<T, any>
): FieldSpecification {
  return getSchemaRuntime(schema, path)
}

export function joinPath<T extends Schema>(path: SchemaPaths<T, 'All'>) {
  return path.reduce((acc: string, p) => {
    let suffix = String(p)
    if (acc.length !== 0 && typeof p === 'string') {
      suffix = `.${suffix}`
    } else if (typeof p === 'number') {
      suffix = `[${suffix}]`
    }
    return `${acc}${suffix}`
  }, '')
}

export type SchemaBinder<T extends Schema> = (...path: SchemaPaths<T, 'All'>) => SchemaBinding

/**
 * Schema binding for form elements. Sets field constraints, client-side validation rules, hints, and classes.
 */
export type SchemaBinding = {
  hint?: string
  min?: number
  max?: number
  minLength?: number
  maxLength?: number
  rules: ((value: any) => true | string)[]
  class?: string | Record<string, boolean> | (string | Record<string, boolean>)[]
}

export function patternRule(pattern: string, errMessage = 'Invalid format') {
  const regex = new RegExp(pattern)
  return (value: string) => {
    return !value || regex.test(value) || errMessage
  }
}

export type ValidationError = ErrorDetail & { path: Array<string | number> }

export type PathPrefix<T extends Schema> = Exclude<
  SchemaPaths<T, 'All'>,
  SchemaPaths<T, 'Terminal'>
>

export type PathComplement<
  T extends Schema,
  Pref extends PathPrefix<T> | undefined,
  Paths extends SchemaPaths<T, 'Terminal'> = SchemaPaths<T, 'Terminal'>
> = undefined extends Pref ? Paths : Paths extends [...NonNullable<Pref>, ...infer R] ? R : never

export type PathJoin<
  T extends Schema,
  P extends PathPrefix<T>,
  Complement extends PathComplement<T, P> = PathComplement<T, P>
> =
  [...P, ...Complement] extends SchemaPaths<T, 'All'>
    ? [...P, ...Complement] & SchemaPaths<T, 'All'>
    : never

function joinPathPrefix<T extends Schema, Prefix extends PathPrefix<T> = PathPrefix<T>>(
  prefix: Prefix,
  path: PathComplement<T, Prefix>
): PathJoin<T, Prefix> {
  return [...prefix, ...path] as PathJoin<T, Prefix>
}

/**
 * Derives validation rules from an OpenAPI schema.
 * Compatible with Vuetify form validation API.
 */
function makeRules({ schema: s, required }: FieldSpecification) {
  const rules: Rule[] = []
  if (required)
    rules.push((value: any) => (!!value || value === 0 ? true : 'This field is required'))

  // Length validation
  if (s?.minLength !== undefined) {
    rules.push(
      (value: string | Array<unknown>) =>
        (value?.length ?? 0) >= s.minLength! ||
        `At least ${s.minLength!} ${s.type == 'string' ? 'character(s)' : 'element(s)'} required`
    )
  }
  if (s?.maxLength !== undefined) {
    rules.push(
      (value: string | Array<unknown>) =>
        (value?.length ?? 0) <= s.maxLength! ||
        `At most ${s.maxLength!} ${s.type == 'string' ? 'character(s)' : 'element(s)'} accepted`
    )
  }

  // Numbers
  if (s?.type == 'number' || s?.type == 'integer') {
    rules.push((value?: string | number) => {
      if (value === undefined || value === null || value === '') return true
      value = Number(value)
      if (s.type == 'integer') return Number.isInteger(value) || `Must be an integer number`
      if (s.type == 'number' && s.format == 'float')
        return Number.isFinite(value) || `Must be a decimal number`
      return Number.isFinite(value) || `Must be a number`
    })
  }
  if (s?.maximum !== undefined) {
    rules.push(
      (value: number) =>
        value === undefined ||
        value === null ||
        value <= s.maximum! ||
        `Maximum value is ${s.maximum!}`
    )
  }
  if (s?.minimum !== undefined) {
    rules.push(
      (value: number) =>
        value === undefined ||
        value === null ||
        value >= s.minimum! ||
        `Minimum value is ${s.minimum!}`
    )
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
    case 'email':
      const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
      rules.push((value: string) => {
        return !value || emailRegex.test(value) || 'Invalid email format'
      })
      break
    // Custom
    case 'country-code':
      rules.push(
        (value: string) =>
          !value || useCountries().findCountry(value) !== undefined || `Invalid country code`
      )
      break

    default:
      break
  }

  return rules
}

/**
 * Bind OpenAPI schema to form elements, providing validation rules.
 * @param schema The OpenAPI schema to bind
 * @returns An object with methods for binding schema to form elements and validating values
 */
export function useSchemaBinding<T extends Schema>(schema: T) {
  /**
   * Traverse schema, collecting all property paths.
   * Used for explicit validation calls using `validateAll`.
   */
  const allPaths = collectPaths(schema)

  /**
   * Validates a value against the OpenAPI schema.
   * In most cases you should use `field` instead, which also binds error messages to form fields.
   * This is useful for validating values outside of form elements.
   */
  function validate(...path: SchemaPaths<T, 'Terminal'>) {
    const spec = getSchema(schema, ...path)
    const rules = makeRules(spec)
    return (value: any) => {
      return rules.reduce<true | string>((acc, rule) => {
        if (acc !== true) return acc
        return rule(value)
      }, true)
    }
  }

  function isArrayPath(path: CollectedPath<T> | SchemaPaths<T>) {
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
          },
          {
            value: v,
            errors: new Array<ValidationError>()
          }
        )
        return errors
      }

      // Get the value at the path
      const value = path.reduce((acc, p) => (acc as Record<string, any>)[p], v)
      const valid = validate(...(path as unknown as SchemaPaths<T, 'Terminal'>))(value)
      return valid !== true
        ? [
            {
              location: path.join('.'),
              message: valid,
              value,
              path: (prefix ?? []).concat(path)
            }
          ]
        : []
    })
  }

  /**
   * Generates field properties from an OpenAPI schema.
   * Compatible with Vuetify form elements.
   */
  function fieldProps(...path: SchemaPaths<T, 'All'>): SchemaBinding {
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
      class: { required: spec.required }
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
  function bindSchema(...path: SchemaPaths<T, 'All'>): SchemaBinding {
    return fieldProps(...path)
  }

  function withPrefix(...prefix: PathPrefix<T>) {
    const fieldBinder = (...path: PathComplement<T, PathPrefix<T>>): SchemaBinding => {
      const fullPath = joinPathPrefix(prefix, path)
      return fieldProps(...fullPath)
    }
    return {
      schema: fieldBinder
    }
  }

  return {
    /**
     * Bind client-side constraints and validation rules.
     */
    schema: bindSchema,
    withPrefix,
    validate,
    validateAll,
    paths: collectPaths(schema)
  }
}

export function addRules<T extends { rules: Rule[] }>({ rules, ...rest }: T, ...add: Rule[]) {
  return {
    ...rest,
    rules: rules.concat(add)
  }
}

export function useSchemaErrors<T extends Schema>(
  _schema: T,
  errors: Ref<AppError | undefined | null>
) {
  const indexedErrors = computed<IndexedValidationErrors>(() => {
    const indexed: IndexedValidationErrors = { rest: [] }
    if (errors.value?.errors) {
      errors.value.errors.forEach((e) => {
        if (e.location === undefined) indexed.rest.push(e.message ?? 'Invalid value')
        else if (e.location.startsWith('body.')) {
          const path = e.location.slice('body.'.length)
          indexed[path] ??= []
          indexed[path].push(e.message ?? 'Invalid value')
        }
      })
    }
    return indexed
  })

  function bindErrors(...path: SchemaPaths<T, 'All'>): string[] | undefined {
    return indexedErrors.value[joinPath(path)]
  }

  return { errors: bindErrors }
}
