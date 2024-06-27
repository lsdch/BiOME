import { ErrorDetail } from "@/api"
import * as Schemas from "@/api/schemas.gen"
import { useCountries } from "@/stores/countries"
import { List, Union } from "ts-toolbelt"


type SchemaModule = typeof import("@/api/schemas.gen")

// Adapted from openapi-ts src
export type Schema = Readonly<{
  $ref?: string,
  additionalProperties?: (boolean | Schema)
  allOf?: Readonly<Schema[]>
  anyOf?: Readonly<Schema[]>
  const?: string | number | boolean | null
  default?: unknown
  deprecated?: boolean
  description?: string
  enum?: Readonly<(string | number)[]>
  example?: unknown
  exclusiveMaximum?: boolean
  exclusiveMinimum?: boolean
  format?: string
  items?: Schema
  maximum?: number
  maxItems?: number
  maxLength?: number
  maxProperties?: number
  minimum?: number
  minItems?: number
  minLength?: number
  minProperties?: number
  multipleOf?: number
  not?: Readonly<Schema[]>
  nullable?: boolean
  oneOf?: Readonly<Schema[]>
  pattern?: string
  properties?: Readonly<Record<string, Schema>>
  readOnly?: boolean
  required?: Readonly<string[]>
  title?: string
  type?: string | Readonly<string[]>
  uniqueItems?: boolean
  writeOnly?: boolean
}>

export type SchemaProperties = Readonly<Record<string, Schema>>
export type SchemaWithProperties<P> = Schema & Readonly<{ type: "object", properties: P }>


/**
 * All property paths in an OpenAPI schema
 */
export type SchemaPaths<T extends Schema | undefined, Terminal extends "Terminal" | "All" = "All"> =
  T extends { properties: Record<string, Schema> } ? {
    [K in keyof T['properties']]-?: (Terminal extends true ? never : [K]) | [K, ...SchemaPaths<T['properties'][K]>]
  }[keyof T['properties']]
  : T extends { items: Schema }
  ? ([number, ...SchemaPaths<T['items']>])
  : T extends { $ref: `#/components/schemas/${infer Z}` }
  ? (`$${Z}` extends keyof SchemaModule ? SchemaPaths<SchemaModule[`$${Z}`]> : [])
  : []

export type StringPath<P extends (string | number)[]> =
  P extends [] ? '' :
  P extends [number] ? `[${P[number]}]` :
  P extends [string] ? `${P[number]}` :
  P extends [number, ...infer K] ? (K extends (string | number)[] ? `[${P[0]}]${K[0] extends string ? '.' : ''}${StringPath<K>}` : never) :
  P extends [string, ...infer K] ? (K extends (string | number)[] ? `${P[0]}${K[0] extends string ? '.' : ''}${StringPath<K>}` : never) :
  never


function paths(s: Schema): (string | '*')[][] {
  if (s.properties) {
    return Object.entries<Schema>(s.properties).reduce<(string | '*')[][]>((acc, [prop, schema]) => {
      if (schema.$ref) {
        const key = `$${schema.$ref.split('/').at(-1)}` as keyof SchemaModule
        // console.log(key, paths(Schemas[key]))
        const p = paths(Schemas[key]).map(p => [prop, ...p])
        return acc.concat(p.length ? p : [[prop]])
      }
      if (schema.properties)
        return acc.concat(paths(schema).map(p => [prop, ...p]))
      if (schema.items)
        return acc.concat(paths(schema).map(p => [prop, '*', ...p]))
      return acc.concat([[prop]])
    }, [])
  }
  return []
}


export type FieldSpecification = { schema: Schema | undefined, required: boolean }

export function getSchema<T extends Schema>(schema: T | undefined, ...path: SchemaPaths<typeof schema, "Terminal" | "All">): FieldSpecification {
  if (schema === undefined)
    return { schema: undefined, required: false }
  if (schema.$ref !== undefined) {
    const refName = `$${schema.$ref.split('/').at(-1)}` as keyof SchemaModule
    const target = Schemas[refName] as Schema
    const p = path as SchemaPaths<typeof target>
    return getSchema<typeof target>(target, ...p)
  }
  const [fragment, ...rest] = path
  if (rest.length == 0) {
    if (typeof fragment === "string")
      return {
        schema: schema.properties?.[fragment],
        required: schema.required?.includes(fragment) ?? false
      }
    else if (typeof fragment === "number" && schema.items !== undefined)
      return getSchema(schema.items)
  }
  if (typeof fragment === "string")
    return getSchema(schema.properties?.[fragment], ...rest)
  else if (typeof fragment === "number" && schema.items !== undefined)
    return getSchema(schema.items, ...rest)
  else return { schema: undefined, required: false }
}


export type SchemaBinding = {
  hint?: string,
  min?: number,
  max?: number,
  minLength?: number,
  maxLength?: number,
  rules: ((value: any) => true | string)[]
}

export function useSchema<T extends Schema>(schema: T) {

  type Rule = ((v: any) => true | string)

  function makeRules({ schema: s, required }: FieldSpecification) {
    const rules: Rule[] = []
    if (required) rules.push((value: any) => value || value === 0 ? true : "This field is required")

    // Strings
    if (s?.minLength !== undefined) {
      rules.push((value: string) => (value?.length ?? 0) >= (s.minLength!) ? true : `At least ${s.minLength!} characters required`)
    }
    if (s?.maxLength !== undefined) {
      rules.push((value: string) => (value?.length ?? 0) <= (s.maxLength!) ? true : `At most ${s.maxLength!} characters accepted`)
    }

    // Numbers
    if ((s?.type == 'number' || s?.type == "integer")) {
      rules.push((value?: number) => {
        if (value === undefined) return true
        if (s.type == "integer") return Number.isInteger(value) || `Must be an integer number`
        if (s.type == "number" && s.format == "float") return (Number.isFinite(value) && !Number.isInteger(value)) || `Must be a decimal number`
        return Number.isFinite(value) || `Must be a number`
      })
    }
    if (s?.maximum !== undefined) {
      rules.push((value: number) => (value <= s.maximum!) || `Maximum value is ${s.maximum!}`)
    }
    if (s?.minimum !== undefined) {
      rules.push((value: number) => (value >= s.minimum!) || `Minimum value is ${s.minimum!}`)
    }

    // Custom
    if (s?.format == "country-code") {
      rules.push((value: string) => useCountries().findCountry(value) !== undefined || `Invalid country code`)
    }

    return rules
  }

  function validate(...path: SchemaPaths<typeof schema>) {
    const spec = getSchema(schema, ...path)
    const rules = makeRules(spec)
    return (value: any) => {
      return rules.reduce<true | string>((acc, rule) => {
        if (acc !== true) return acc
        return rule(value)
      }, true)
    }
  }


  const allPaths = paths(schema) as unknown as Union.ListOf<List.Replace<SchemaPaths<T, "All">, number, '*'>>

  function validateAll(v: Record<string, any>) {
    return allPaths.flatMap<(ErrorDetail & { path: string[] })>((path: string[]): (ErrorDetail & { path: string[] })[] => {
      if (path.includes('*')) {
        return []
      }
      const value = path.reduce((acc, p) => acc[p], v)
      const valid = validate(...path as (typeof allPaths)[number])(value)
      return valid !== true ? [{ location: path.join('.'), message: valid, value, path }] : []
    })
  }

  function bindSchema(...path: SchemaPaths<typeof schema, "Terminal">): SchemaBinding {
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
    }
  }

  return { schema: bindSchema, validate, paths: paths(schema), validateAll }
}