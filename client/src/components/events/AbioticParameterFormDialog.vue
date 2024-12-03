<template>
  <FormDialog v-model="dialog" :title="`${mode} abiotic parameter`" :loading @submit="submit">
    <v-container fluid>
      <v-row>
        <v-col>
          <v-text-field label="Label" v-model="model.label" v-bind="field('label')" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field label="Code" v-model="model.code" v-bind="field('code')" />
        </v-col>
        <v-col>
          <v-text-field label="Unit" v-model="model.unit" v-bind="field('unit')" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea
            label="Description"
            v-model="model.description"
            v-bind="field('description')"
          />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import {
  $AbioticParameterInput,
  $AbioticParameterUpdate,
  AbioticParameter,
  AbioticParameterInput,
  AbioticParameterUpdate,
  SamplingService
} from '@/api'
import { FormEmits, FormProps, useForm, useSchema } from '@/components/toolkit/forms/form'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { reactiveComputed, useToggle } from '@vueuse/core'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<AbioticParameter>>()
const emit = defineEmits<FormEmits<AbioticParameter>>()

const initial: AbioticParameterInput = {
  label: '',
  code: '',
  unit: ''
}

const { model, mode, makeRequest } = useForm(props, {
  initial,
  updateTransformer({ id, $schema, meta, ...rest }): AbioticParameterUpdate {
    return rest
  }
})

const { field, errorHandler } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? $AbioticParameterInput : $AbioticParameterUpdate)
)

const [loading, toggleLoading] = useToggle(false)

async function submit() {
  toggleLoading(true)
  return await makeRequest({
    create: SamplingService.createAbioticParameter,
    edit: ({ code }, body) => SamplingService.updateAbioticParameter({ path: { code }, body })
  })
    .then(errorHandler)
    .then((item) => emit('success', item))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped lang="scss"></style>
