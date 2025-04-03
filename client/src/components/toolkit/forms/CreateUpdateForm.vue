<template>
  <slot
    name="default"
    :mode
    :model
    :field
    :loading="activeMutation.isPending"
    :submit
    :setModel
    v-bind="$attrs"
  />
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
    InputRequestData extends RequestData<ItemInput>,
    UpdateRequestData extends RequestData<ItemUpdate> & { path: ItemID },
    InputModel = ItemInput,
    UpdateModel = ItemUpdate
  "
>
import { ErrorModel } from '@/api'
import { FormCreateMutation, FormUpdateMutation, RequestData } from '@/functions/mutations'
import { useMutation } from '@tanstack/vue-query'
import { reactiveComputed } from '@vueuse/core'
import { StatusCodes } from 'http-status-codes'
import { computed, ref, watch } from 'vue'
import { Schema, useSchema } from '../../../composables/schema'

export type Mode = 'Create' | 'Edit'

const { create, update, local } = defineProps<{
  create: FormCreateMutation<Item, ItemInput, InputModel, InputSchema, InputRequestData>
  update: FormUpdateMutation<Item, ItemUpdate, UpdateModel, UpdateSchema, ItemID, UpdateRequestData>
  local?: true
}>()

const emit = defineEmits<{
  (evt: 'success', item: Item): void
  (evt: 'created', item: Item): void
  (evt: 'updated', item: Item): void
  (evt: 'save', item: ItemInput): void
  (evt: 'error', item: ErrorModel): void
}>()

const item = defineModel<Item>()
const model = ref<InputModel | UpdateModel>(initModel(item.value))
const mode = computed<Mode>(() => (item.value ? 'Edit' : 'Create'))
watch(item, (item) => (model.value = initModel(item)), { immediate: true })

const { bindSchema: field, dispatchErrors } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? create.schema : update.schema)
)

const createMutation = useMutation({
  ...create.mutation,
  onSuccess(data) {
    emit('success', data)
    emit('created', data)
  },
  onError
})
const updateMutation = useMutation({
  ...update.mutation,
  onSuccess(data) {
    emit('success', data)
    emit('updated', data)
  },
  onError
})

const activeMutation = computed(() => (mode.value === 'Create' ? createMutation : updateMutation))

async function submit() {
  if (local) {
    emit('save', create.requestData?.(model.value)?.body ?? (model.value as ItemInput))
    return
  }
  if (mode.value === 'Create')
    return await createMutation.mutateAsync(create.requestData(model.value))
  else {
    return await updateMutation.mutateAsync(update.requestData(item.value!, model.value))
  }
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
  return item ? update.itemToModel(item) : create.initial
}

function reset() {
  initModel(item.value)
}
</script>

<style scoped lang="scss"></style>
