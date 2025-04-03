<template>
  <HabitatGroupFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="title(mode)"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
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
import { defineFormCreate, defineFormUpdate, Mode, useMutationForm } from '@/functions/mutations'
import { HabitatModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import HabitatGroupFormDialog from './HabitatGroupFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<HabitatGroup>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createHabitatGroupMutation(), {
  initial: HabitatModel.initialModel,
  schema: $HabitatGroupInput,
  requestData: (model) => ({ body: HabitatModel.toCreateRequestBody(model) })
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

const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    dialog.value = false
    feedback({
      type: 'success',
      message: mode === 'Create' ? `Habitat group created` : `Habitat group updated`
    })
  }
})

function title(mode: Mode) {
  return mode == 'Create' ? 'Create habitat group' : `Edit habitats: ${item.value!.label}`
}
</script>

<style scoped lang="scss"></style>
