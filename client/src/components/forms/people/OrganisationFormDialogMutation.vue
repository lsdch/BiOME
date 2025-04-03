<template>
  <OrganisationFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} site`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  ></OrganisationFormDialog>
</template>

<script setup lang="ts">
import { $OrganisationInput, $OrganisationUpdate, Organisation } from '@/api'
import {
  createOrganisationMutation,
  updateOrganisationMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { OrganisationModel } from '@/models'
import OrganisationFormDialog from './OrganisationFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Organisation>()

defineProps<FormDialogProps>()

const create = defineFormCreate(createOrganisationMutation(), {
  initial: OrganisationModel.initialModel,
  schema: $OrganisationInput
})

const update = defineFormUpdate(updateOrganisationMutation(), {
  schema: $OrganisationUpdate,
  itemToModel: OrganisationModel.fromOrganisation,
  requestData: ({ code }) => ({ path: { code } })
})

const { mode, model, errors, submit, activeMutation } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    dialog.value = false
  }
})
</script>

<style scoped></style>
