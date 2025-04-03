<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} data source`"
    @submit="emit('submit', model)"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
    <v-container fluid>
      <v-row>
        <v-col>
          <v-text-field label="Label" v-model.trim="model.label" v-bind="schema('label')" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <FTextField
            label="Code"
            v-model.trim="model.code"
            v-bind="schema('code')"
            class="input-font-monospace"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field
            label="URL"
            v-model.trim="model.url"
            v-bind="schema('url')"
            class="input-font-monospace"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field
            label="Link template"
            v-model.trim="model.link_template"
            v-bind="schema('link_template')"
            class="input-font-monospace"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea
            label="Description"
            v-model.trim="model.description"
            v-bind="schema('description')"
          />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { $DataSourceInput, $DataSourceUpdate } from '@/api'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { useSchema } from '@/composables/schema'
import { FormProps } from '@/functions/mutations'
import { DataSourceModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'
import FTextField from '../toolkit/forms/FTextField'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<DataSourceModel.DataSourceFormModel>({
  default: DataSourceModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: DataSourceModel.DataSourceFormModel | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $DataSourceInput : $DataSourceUpdate))
</script>

<style scoped lang="scss"></style>
