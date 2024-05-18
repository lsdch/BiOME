import { ApiError, ErrorModel } from "@/api"

import { Ref, reactive, ref, toValue, watch } from "vue"
import { Schema, SchemaBinding, useSchema, type SchemaPaths } from "./schema"

export * from "./FormDialog.vue"
export * from "./schema"

export type Mode = 'Create' | 'Edit'

export type FormProps<ItemType> = {
  edit?: ItemType
}

export type FormEmits<ItemType> = {
  (evt: "success", item: ItemType): void
}

function joinPath<T extends Schema>(path: SchemaPaths<T>) {
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
export type FieldBinding = ErrorBinding & SchemaBinding


/**
 * Provides utility functions for client-side input validation,
 * and display of server-side validation errors
 * @param schema The OpenAPI schema for the form data model
 */
export function useForm<
  ItemType extends { [k: string]: unknown },
  ItemInputType extends Partial<Record<keyof ItemType, unknown>>,
  T extends Schema,
>(
  props: FormProps<ItemType>,
  schema: T,
  dataModel: {
    initial: ItemInputType,
    transformers?: Partial<{
      [k in keyof (ItemType | ItemInputType)]: (v: ItemType[k]) => ItemInputType[k]
    }>
  }
) {

  const model = ref(dataModel.initial) as Ref<ItemInputType>

  const { schema: bindSchema } = useSchema<T>(schema)

  const loading = ref(false)

  watch(() => props.edit, (item) => {
    if (item === undefined) {
      model.value = dataModel.initial
    } else {
      const it = toValue(item)
      model.value = Object.fromEntries(
        Object.keys(toValue(model)).map(
          (k: keyof (ItemType | ItemInputType)) =>
            [k, dataModel.transformers?.[k]?.(it[k]) ?? it[k]]
        )
      ) as ItemInputType
    }
  })

  /**
   * Input validation errors indexed by their object path in the API request body
   */
  const errors = ref<Record<string, string[]>>({})

  /**
   * Collects error messages indexed by their object path in an API request body,
   * so that they can be consumed by `bindErrors` or `field`.
   * @param error An instance of `ApiError` returned from an API request.
   */
  function errorHandler({ body }: { body: ErrorModel } & ApiError) {
    body.errors?.forEach(({ location, message }) => {
      if (location === undefined || message === undefined) return
      if (location.startsWith('body.')) {
        const loc = location.replace('body.', '')
        errors.value[loc].push(message)
      }
    })
  }

  /**
   * Binds remote error messages to an input form element.
   * Errors must be caught using `errorHandler` function.
   *
   * @param path The object property path for the field
   */
  function bindErrors(...path: SchemaPaths<T>): ErrorBinding {
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
   * @example `<v-text-field v-model="model.someArray[0].someProperty" v-bind="field('someArray', 0, 'someProperty')"/>
   * @param path The object property path for the field
   * @returns Field bindings to be passed to form element using `v-bind`
   */
  function field(...path: SchemaPaths<T>): FieldBinding {
    return {
      ...bindSchema(...path),
      ...bindErrors(...path),
    }
  }


  return { errors, loading, errorHandler, bindErrors, field, bindSchema, model }
}
