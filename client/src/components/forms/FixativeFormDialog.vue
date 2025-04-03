<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} fixative`"
    @submit="emit('submit', model)"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
    <v-container fluid>
      <v-row>
        <v-col>
          <v-text-field label="Label" v-model="model.label" v-bind="schema('label')" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field label="Code" v-model="model.code" v-bind="schema('code')" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea
            label="Description"
            v-model="model.description"
            v-bind="schema('description')"
          />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { $FixativeInput, $FixativeUpdate } from '@/api'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { useSchema } from '@/composables/schema'
import { FormProps } from '@/functions/mutations'
import { FixativeModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<FixativeModel.FixativeFormModel>({
  default: FixativeModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: FixativeModel.FixativeFormModel | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $FixativeInput : $FixativeUpdate))
</script>

<style scoped lang="scss"></style>
