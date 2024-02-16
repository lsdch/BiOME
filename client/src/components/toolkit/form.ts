import { ApiError, CancelablePromise, InputValidationError } from "@/api"
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

  type ErrorMsgs<ItemInputType> = { [k in keyof ItemInputType]: any }

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
      .then((item: ItemType) => {
        emit("success", item)
      })
      ?.catch((error: ApiError) => {
        switch (error.status) {
          case HttpStatusCode.BadRequest:
            if (typeof error.body === 'object') {
              console.log(error.body)
              errors.value = error.body
            }
            break

          default:
            break
        }
      })
      .finally(() => {
        loading.value = false
      })
  }

  return { errors, loading, submit, errorMsgs }
}