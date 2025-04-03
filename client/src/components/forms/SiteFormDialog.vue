<template>
  <FormDialog
    v-bind="props"
    @submit="emit('submit', model)"
    v-model="dialog"
    :title="title ?? `${mode} site`"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData" />
    </template>
    <SiteForm v-model="model" :mode :errors />
  </FormDialog>
</template>

<script setup lang="ts">
import { FormProps } from '@/functions/mutations'
import { SiteModel } from '@/models'
import FormDialog, { FormDialogProps } from '../toolkit/forms/FormDialog.vue'
import SiteForm from './SiteForm.vue'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<SiteModel.SiteFormModel>({
  default: SiteModel.initialModel
})

const props = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: SiteModel.SiteFormModel | undefined]
}>()
</script>

<style scoped lang="scss"></style>
