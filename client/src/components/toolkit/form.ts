import { ApiError, CancelablePromise, InputValidationError } from "@/api"
import { HttpStatusCode } from "axios"

import { computed } from "vue"
import { Ref, ref } from "vue"

export type Mode = 'Create' | 'Edit'

export type Props<ItemType> = {
  edit?: ItemType
  onSuccess?: (mode: Mode, inst: ItemType) => any
}

export type Emits<ItemType> = {
  success: [mode: Mode, item: ItemType]
}

export function useForm<ItemInputType, ItemType>
  (props: Props<ItemType>,
    emit: (evt: "success", mode: Mode, item: ItemType) => void,
    submitRequest: () => CancelablePromise<ItemType> | Promise<ItemType>) {
  const loading = ref(false)
  const mode = computed((): Mode => (props.edit ? 'Edit' : 'Create'))

  const errors: Ref<Partial<{
    [Property in keyof ItemInputType]: InputValidationError[]
  }>> = ref({})

  function submit() {
    loading.value = true
    submitRequest()
      .then((item: ItemType) => emit("success", mode.value, item))
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

  return { errors, mode, loading, submit }
}