<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} gene`"
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
          <v-text-field
            label="Code"
            v-model.trim="model.code"
            v-bind="schema('code')"
            class="input-font-monospace"
          />
        </v-col>
        <v-col>
          <v-switch
            label="MOTU delimiter"
            v-model="model.is_MOTU_delimiter"
            color="primary"
            v-bind="schema('is_MOTU_delimiter')"
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
import { $GeneInput, $GeneUpdate } from '@/api'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { useSchema } from '@/composables/schema'
import { FormProps } from '@/functions/mutations'
import { GeneModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<GeneModel.GeneFormModel>({
  default: GeneModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: GeneModel.GeneFormModel | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $GeneInput : $GeneUpdate))
</script>

<style scoped lang="scss"></style>
