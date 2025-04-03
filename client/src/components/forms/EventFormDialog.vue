<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} site`"
    @submit="emit('submit', model)"
  >
    <template #subtitle>
      <v-chip label :text="site.code" class="font-monospace" prepend-icon="mdi-map-marker" />
    </template>
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>

    <EventForm v-model="model" :site :mode :errors />
  </FormDialog>
</template>

<script setup lang="ts">
import { SiteItem } from '@/api'
import { FormProps } from '@/functions/mutations'
import { EventModel } from '@/models'
import { SiteFormModel } from '@/models/site'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import EventForm from '@/components/forms/EventForm.vue'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<EventModel.EventModel>({
  default: EventModel.initialModel
})
const emit = defineEmits<{
  submit: [model: EventModel.EventModel | undefined]
}>()
const props = defineProps<{ site: SiteItem | SiteFormModel } & FormProps & FormDialogProps>()
</script>

<style scoped lang="scss"></style>
