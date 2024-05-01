import { ApiError, CancelablePromise } from "@/api"
import { HttpStatusCode } from "axios"
import { computed } from "vue"

import { Ref, ref } from "vue"

export type Mode = 'Create' | 'Edit'

export type Props<ItemType> = {
  edit?: ItemType
}

export type Emits<ItemType> = {
  (evt: "success", item: ItemType): void
}

export type ValidationErrors<ItemInputType> = Partial<
  Record<keyof ItemInputType, InputValidationError[]>
>

export type ErrorMsgs<ItemInputType> = { [k in keyof ItemInputType]: any }

export function useForm<ItemInputType extends Record<string | symbol, any>, ItemType>(
  props: Props<ItemType>,
  emit: Emits<ItemType>,
  submitRequest: () => CancelablePromise<ItemType> | Promise<ItemType>
) {
  const loading = ref(false)

  /**
   * Input validation errors indexed by field name
   */
  const errors: Ref<ValidationErrors<ItemInputType>> = ref({})

  /**
   * A proxy to validation errors that allows direct access to error messages
   */
  const errorMsgs = computed((): ErrorMsgs<ItemInputType> => {
    return Object.fromEntries(
      Object.entries(errors.value).map(
        ([key, val]) => {
          return [key, val?.map(({ message }) => message)]
        }
      )
    ) as ErrorMsgs<ItemInputType>
  })



  function submit() {
    loading.value = true
    submitRequest()
      .then((item: ItemType) => { emit("success", item) })
      ?.catch((error: ApiError) => {
        switch (error.status) {
          case HttpStatusCode.BadRequest:
            if (typeof error.body === 'object') {
              errors.value = error.body
            }
            break

          default:
            break
        }
      })
      .finally(() => { loading.value = false })
  }

  return { errors, loading, submit, errorMsgs }
}



import { ValidationRuleWithoutParams } from '@vuelidate/core'

export function inlineRule(rule: ValidationRuleWithoutParams) {
  return (value: any) => {
    return rule.$validator(value, undefined, undefined) ? true : rule.$message?.toString() ?? ''
  }
}

export function inlineRules(rules: ValidationRuleWithoutParams[]) {
  return rules.map(inlineRule)
}


// Adapted from openapi-ts src
export type Schema = Readonly<{
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
export function useSchema<P extends SchemaProperties>(schema: SchemaWithProperties<P>) {
  function inputProps(key: keyof P) {
    const k = key as unknown as string
    const s = schema.properties?.[k]

    return {
      hint: s?.description,
      min: s?.minimum,
      max: s?.maximum,
      rules: schema.required?.includes(k)
        ? [
          (value: any) => value || value === 0 ? true : "This field is required"
        ]
        : undefined
    }
  }

  return { schema: inputProps }
}