<template>
  <slot name="default" :mode :model :field :loading="activeMutation.isPending" :submit :setModel />
</template>

<script
  setup
  lang="ts"
  generic="
    Item,
    ItemInput,
    ItemUpdate,
    ItemID,
    InputSchema extends Schema,
    UpdateSchema extends Schema,
    InputModel = ItemInput,
    UpdateModel = ItemUpdate
  "
>
import { ErrorModel } from '@/api'
import { useMutation, UseMutationOptions } from '@tanstack/vue-query'
import { reactiveComputed } from '@vueuse/core'
import { StatusCodes } from 'http-status-codes'
import { computed, ref, watch } from 'vue'
import { Schema, useSchema } from './schema'
import { Equal } from 'node_modules/@tanstack/vue-query/build/modern/types'

export type Mode = 'Create' | 'Edit'

export type FormCreateMutation<
  Item = any,
  ItemInput = any,
  InputModel = ItemInput,
  InputSchema = Schema
> = {
  mutation: (options?: {
    body: ItemInput
  }) => UseMutationOptions<Item, ErrorModel, { body: ItemInput }, any>
  schema: InputSchema
} & (Equal<InputModel, ItemInput> extends true
  ? {}
  : { transformer: (model: InputModel) => ItemInput })

export type FormUpdateMutation<
  Item = any,
  ItemUpdate = any,
  UpdateModel = ItemUpdate,
  UpdateSchema = Schema,
  ItemID = any
> = {
  mutation: (options?: { path: ItemID; body: ItemUpdate }) => UseMutationOptions<
    Item,
    ErrorModel,
    {
      path: ItemID
      body: ItemUpdate
    },
    any
  >
  schema: UpdateSchema
  itemID(item: Item): ItemID
} & (Equal<UpdateModel, ItemUpdate> extends true
  ? {}
  : { transformer: (item: UpdateModel) => ItemUpdate })

const { initial, updateTransformer, create, update } = defineProps<{
  initial: InputModel
  updateTransformer: (item: Item) => UpdateModel
  create: FormCreateMutation<Item, ItemInput, InputModel, InputSchema>
  update: FormUpdateMutation<Item, ItemUpdate, UpdateModel, UpdateSchema, ItemID>
}>()

const emit = defineEmits<{
  (evt: 'success', item: Item): void
  (evt: 'created', item: Item): void
  (evt: 'updated', item: Item): void
  (evt: 'error', item: ErrorModel): void
  (evt: 'loading', loading: boolean): void
}>()

const item = defineModel<Item>()
const model = ref<InputModel | UpdateModel>(initModel(item.value))
const mode = computed<Mode>(() => (item.value ? 'Edit' : 'Create'))
watch(item, (item) => (model.value = initModel(item)), { immediate: true })

const { field, dispatchErrors } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? create.schema : update.schema)
)

const createMutation = useMutation({
  ...create.mutation(),
  onSuccess(data) {
    emit('success', data)
    emit('created', data)
  },
  onError
})
const updateMutation = useMutation({
  ...update.mutation(),
  onSuccess(data) {
    emit('success', data)
    emit('updated', data)
  },
  onError
})

const activeMutation = computed(() => (mode.value === 'Create' ? createMutation : updateMutation))
watch(activeMutation.value.isPending, (pending) => emit('loading', pending))

async function submit() {
  if (mode.value === 'Create')
    return await createMutation.mutateAsync({ body: model.value as ItemInput })
  else
    return await updateMutation.mutateAsync({
      path: update.itemID(item.value!),
      body: model.value as ItemUpdate
    })
}

function setModel(newModel: InputModel | UpdateModel) {
  model.value = newModel
}

defineExpose({ submit, reset })

function onError(err: ErrorModel) {
  if (err.status && err.status !== StatusCodes.UNPROCESSABLE_ENTITY) {
    dispatchErrors(err)
  }
  emit('error', err)
}

function initModel(item?: Item) {
  return item ? updateTransformer(item) : initial
}
function reset() {
  initModel(item.value)
}
</script>

<style scoped lang="scss"></style>
