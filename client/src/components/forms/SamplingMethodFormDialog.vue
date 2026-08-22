<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} sampling method`"
    @submit="emit('submit', model)"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
    <v-container fluid>
      <v-row>
        <v-col>
          <v-text-field label="Label" v-model="model.name" v-bind="schema('name')" />
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
import { $SamplingMethodInput, $SamplingMethodUpdateParams } from '@/api'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { useSchemaBinding } from '@/composables/schema'
import { FormProps } from '@/lib/mutations'
import { SamplingMethodModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<SamplingMethodModel.SamplingMethodFormModel>({
  default: SamplingMethodModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: SamplingMethodModel.SamplingMethodFormModel | undefined]
}>()

const { schema } = reactiveComputed(() =>
  useSchemaBinding(mode === 'Create' ? $SamplingMethodInput : $SamplingMethodUpdateParams)
)
</script>

<style scoped lang="scss"></style>
