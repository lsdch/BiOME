<template>
  <HabitatGroupFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
    :depends
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </HabitatGroupFormDialog>
</template>

<script setup lang="ts">
import { $HabitatGroupInput, $HabitatGroupUpdate, HabitatGroup } from '@/api'
import {
  createHabitatGroupMutation,
  updateHabitatGroupMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/lib/mutations'
import { HabitatModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import { computed } from 'vue'
import HabitatGroupFormDialog from './HabitatGroupFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<HabitatGroup>('item')

const { depends } = defineProps<FormDialogProps & { depends?: string }>()

const create = defineFormCreate(createHabitatGroupMutation(), {
  initial: HabitatModel.initialModel,
  schema: $HabitatGroupInput,
  requestData: (model) => ({ body: HabitatModel.toCreateRequestBody({ ...model, depends }) })
})

const update = defineFormUpdate(updateHabitatGroupMutation(), {
  itemToModel: HabitatModel.fromHabitatGroup,
  schema: $HabitatGroupUpdate,
  requestData: ({ label }, model) => ({
    path: { label },
    body: HabitatModel.toUpdateRequestBody(model)
  })
})

const { feedback } = useFeedback()
const emit = defineEmits<{
  success: [item: HabitatGroup]
}>()
const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    emit('success', item)
    dialog.value = false
    feedback({
      type: 'success',
      message: mode === 'Create' ? `Habitat group created` : `Habitat group updated`
    })
  }
})

const title = computed(() => {
  if (mode.value === 'Create') return 'Create habitat group'
  return `Edit habitat group: ${item.value!.label}`
})
</script>

<style scoped lang="scss"></style>
