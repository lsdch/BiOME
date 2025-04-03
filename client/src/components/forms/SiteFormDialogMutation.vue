<template>
  <SiteFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} site`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </SiteFormDialog>
</template>

<script setup lang="ts">
import { $SiteInput, $SiteUpdate, Site } from '@/api'
import { createSiteMutation, updateSiteMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { SiteModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import SiteFormDialog from './SiteFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Site>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createSiteMutation(), {
  initial: SiteModel.initialModel,
  schema: $SiteInput,
  requestData: (model) => ({ body: SiteModel.toRequestBody(model) })
})

const update = defineFormUpdate(updateSiteMutation(), {
  itemToModel: SiteModel.fromSite,
  schema: $SiteUpdate,
  requestData: ({ code }, model) => ({
    path: { code },
    body: SiteModel.toRequestBody(model)
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
      message: mode === 'Create' ? `Site created` : `Site updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
