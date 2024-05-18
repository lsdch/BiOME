import * as Schemas from "@/api/schemas.gen"
import { ValidationRuleWithoutParams } from "@vuelidate/core"


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

export type SchemaProperties = Readonly<Record<string, {}>>
export type SchemaWithProperties<P> = Schema & Readonly<{ type: "object", properties: P }>



export type SchemaPaths<T extends Schema | undefined> =
  T extends { properties: Record<string, Schema> } ? {
    [K in keyof T['properties']]-?: [K] | [K, ...SchemaPaths<T['properties'][K]>]
  }[keyof T['properties']]
  : T extends { items: Schema }
  ? ([number, ...SchemaPaths<T['items']>])
  : T extends { $ref: `#/components/schemas/${infer Z}` }
  ? (`$${Z}` extends keyof SchemaModule ? SchemaPaths<SchemaModule[`$${Z}`]> : [])
  : []




export type FieldSpecification = { schema: Schema | undefined, required: boolean }

export function getSchema<T extends Schema>(schema: T | undefined, ...path: SchemaPaths<typeof schema> | []): FieldSpecification {
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

  function bindSchema(...path: SchemaPaths<typeof schema> | []): SchemaBinding {
    const { schema: s, required } = getSchema(schema, ...path)

    const rules = []
    if (required) rules.push((value: any) => value || value === 0 ? true : "This field is required")
    if (s?.minLength !== undefined) {
      rules.push((value: string) => (value?.length ?? 0) >= (s.minLength ?? 0) ? true : `At least ${s.minLength} characters required`)
    }

    return {
      hint: s?.description,
      min: s?.minimum,
      max: s?.maximum,
      minLength: s?.minLength,
      maxLength: s?.maxLength,
      rules,
    }
  }

  return { schema: bindSchema }
}