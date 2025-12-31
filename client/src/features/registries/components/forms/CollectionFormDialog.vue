<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} Collection`"
    @submit="emit('submit', model)"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
    <v-container fluid>
      <v-text-field v-model="model.label" v-bind="schema('label')" label="Label"></v-text-field>
      <v-text-field v-model="model.code" v-bind="schema('code')" label="Code" />
      <v-text-field
        v-model="model.contact"
        v-bind="schema('contact')"
        label="Contact"
      ></v-text-field>
      <v-text-field
        v-model="model.location"
        v-bind="schema('location')"
        label="Location"
      ></v-text-field>
      <v-checkbox
        v-model="model.personal"
        v-bind="schema('personal')"
        label="Personal Collection"
      ></v-checkbox>
      <v-textarea
        v-model="model.description"
        v-bind="schema('description')"
        label="Description"
        rows="3"
      ></v-textarea>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { $CollectionInput, $CollectionUpdate } from '@/api'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { useSchema } from '@/composables/schema'
import { FormProps } from '@/lib/mutations'
import { CollectionModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<CollectionModel.CollectionFormModel>({
  default: CollectionModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: CollectionModel.CollectionFormModel | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $CollectionInput : $CollectionUpdate))
</script>

<style scoped lang="scss"></style>
