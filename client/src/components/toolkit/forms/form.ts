import { ErrorModel } from "@/api"

import { Ref, computed, reactive, ref, toValue, watch } from "vue"
import { Schema, SchemaBinding, useSchema, type SchemaPaths } from "./schema"
import { ResponseBody, handleErrors } from "@/api/responses"

export * from "./FormDialog.vue"
export * from "./schema"

export type Mode = 'Create' | 'Edit'

export type FormProps<ItemType> = {
  edit?: ItemType
}

export type FormEmits<ItemType> = {
  (evt: "success", item: ItemType): void
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
>(
  props: FormProps<ItemType>,
  dataModel: {
    /**
     * Base value for data model
     */
    initial: ItemInputType,
    /**
     * Transform item fields to their representation for update inputs
     */
    transformers?: Partial<{
      [k in keyof (ItemType | ItemInputType)]: (v: ItemType[k]) => ItemInputType[k]
    }>
  }
) {

  const model = ref(dataModel.initial) as Ref<ItemInputType>
  const mode = computed<Mode>(() => props.edit == undefined ? 'Create' : 'Edit')

  watch(() => props.edit, (item) => {
    if (item === undefined) {
      model.value = dataModel.initial
    } else {
      const it = toValue(item)
      model.value = it
      //   Object.fromEntries(
      //   Object.keys(toValue(model)).map(
      //     (k: keyof (ItemType | ItemInputType)) =>
      //       [k, dataModel.transformers?.[k]?.(it[k]) ?? it[k]]
      //   )
      // ) as ItemInputType
    }
  })

  function formTitle(objectName: string) {
    return computed(() => {
      `${mode.value} ${objectName}`
    })
  }


  return { model, mode, formTitle }
}
