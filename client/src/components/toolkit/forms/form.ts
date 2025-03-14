import { ErrorModel } from '@/api'
import { MutationOptions, useMutation, UseMutationOptions } from '@tanstack/vue-query'
import { reactiveComputed } from '@vueuse/core'
import { computed, Ref, ref, toRef, watch } from 'vue'

export * from './schema'

export type Mode = 'Create' | 'Edit'

export type FormProps<Item> = { edit?: Item }

export type FormEmits<ItemType> = {
  (evt: 'success', item: ItemType): void
  (evt: 'created', item: ItemType): void
  (evt: 'updated', item: ItemType): void
}

/**
 * Handles switching from create mode from an empty model
 * to edit mode with initial values from an existing object
 */
export function useForm<
  Item extends object,
  ItemInput,
  ItemUpdate,
  ItemID = unknown,
>(
  props: FormProps<Item>,
  opts: {
    /**
     * Base value for data model
     */
    initial: ItemInput
    /**
     * Defines how to extract initial form values from an existing item
     * so they can be updated
     */
    updateTransformer(item: Item): ItemUpdate

    mutations: {
      create: (options?: { body: ItemInput }) => UseMutationOptions<Item, ErrorModel, { body: ItemInput }, any>,
      update: (options?: { path: ItemID, body: ItemUpdate }) => UseMutationOptions<Item, ErrorModel, {
        path: ItemID, body: ItemUpdate
      }, any>
      itemID(item: Item): ItemID,
      options?: Omit<MutationOptions<Item, ErrorModel, { body: ItemInput | ItemUpdate }>, "mutationFn" | "mutationKey">
    }
  },
) {

  const item = toRef(() => props.edit)
  const model = ref(initModel(item.value)) as Ref<ItemInput | ItemUpdate>
  const mode = computed(() => item.value ? "Edit" : "Create")
  watch(item, (item) => model.value = initModel(item), { immediate: true })


  function initModel(item?: Item) {
    return item
      ? opts.updateTransformer(item)
      : opts.initial
  }


  function reset() {
    initModel(props.edit)
  }


  const createMutation = useMutation({
    ...opts.mutations.create(),
    ...opts.mutations.options
  })
  const updateMutation = useMutation({
    ...opts.mutations.update(),
    ...opts.mutations.options
  })

  const { error, isPending } = reactiveComputed(() => mode.value === "Create" ? createMutation : updateMutation)
  async function makeRequest() {
    if (mode.value === "Create")
      return await createMutation.mutateAsync({ body: model.value as ItemInput })
    else
      return await updateMutation.mutateAsync({
        path: opts.mutations.itemID(item.value!),
        body: model.value as ItemUpdate
      })
  }

  return { mode, model, makeRequest, reset, isPending, error }
}
