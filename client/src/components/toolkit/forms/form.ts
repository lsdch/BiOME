import { ErrorModel } from '@/api'
import { RequestResult } from '@hey-api/client-fetch'
import { computed, Ref, ref, toRef, watch } from 'vue'

export * from './FormDialog.vue'
export * from './schema'

export type Mode = 'Create' | 'Edit'

export type FormProps<Item> = {
  edit?: Item
}

export type FormEmits<ItemType> = {
  (evt: 'success', item: ItemType): void
}

// type T<Item, ItemInput, ItemUpdate> = {
//   mode: ComputedRef<"Create">
//   model: Ref<ItemInput>
//   item?: undefined
// } | {
//   mode: ComputedRef<"Edit">
//   model: Ref<ItemUpdate>
//   item: Item
// }

export type CreateReq<ItemInput, Item> =
  ({ body }: { body: ItemInput }) => RequestResult<Item, ErrorModel, false>
export type UpdateReq<ItemUpdate, Item> =
  (item: Item, model: ItemUpdate) => RequestResult<Item, ErrorModel, false>

/**
 * Handles switching from create mode from an empty model
 * to edit mode with initial values from an existing object
 */
export function useForm<Item extends object, ItemInput, ItemUpdate>(
  props: FormProps<Item>,
  dataModel: {
    /**
     * Base value for data model
     */
    initial: ItemInput
    /**
     * Defines how to extract initial form values from an existing item
     * so they can be updated
     */
    updateTransformer(item: Item): ItemUpdate
  }
) {

  const item = toRef(() => props.edit)

  function initModel(item?: Item) {
    return item
      ? dataModel.updateTransformer(item)
      : dataModel.initial
  }

  const mode = computed(() => item.value ? "Edit" : "Create")
  const model = ref(initModel(item.value)) as Ref<ItemInput | ItemUpdate>

  watch(item, (item) => model.value = initModel(item), { immediate: true })

  function makeRequest({ create, edit }: { create: CreateReq<ItemInput, Item>, edit: UpdateReq<ItemUpdate, Item> }) {
    switch (mode.value) {
      case "Create":
        return create({ body: model.value as ItemInput })
      case "Edit":
        return edit(item.value!, model.value as ItemUpdate)
      default:
        throw `Unexpected form mode value: ${mode.value}`
    }
  }

  return { mode, model, makeRequest }

  // const state = reactiveComputed(() => {
  //   if (editItem.value === undefined) {
  //     return { mode: "Create" as const, model: ref(dataModel.initial) as Ref<ItemInput> }
  //   } else {
  //     return {
  //       mode: "Edit" as const,
  //       model: ref(dataModel.updateTransformer(editItem.value)) as Ref<ItemUpdate>,
  //       item: editItem.value!
  //     }
  //   }
  // })

  // return state

  // const model =
  //   props.edit !== undefined
  //     ? ref(dataModel.updateTransformer(props.edit)) as Ref<ItemUpdate>
  //     : ref(dataModel.initial) as Ref<ItemInput>

  // const mode = computed<Mode>(() => props.edit == undefined ? 'Create' : 'Edit')

  // watch(
  //   () => props.edit,
  //   (item) => {
  //     model.value = (item !== undefined
  //       ? dataModel.updateTransformer(toValue(item))
  //       : dataModel.initial
  //     ) as undefined extends Item ? ItemInput : ItemUpdate
  //   })

  // return { model, mode }
}
