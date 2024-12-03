<template>
  <FormDialog v-model="dialog" :title="`${mode} fixative`" :loading @submit="submit">
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
  $FixativeInput,
  $FixativeUpdate,
  Fixative,
  FixativeInput,
  FixativeUpdate,
  SamplingService
} from '@/api'
import { FormEmits, FormProps, useForm, useSchema } from '@/components/toolkit/forms/form'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { reactiveComputed, useToggle } from '@vueuse/core'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<Fixative>>()
const emit = defineEmits<FormEmits<Fixative>>()

const initial: FixativeInput = {
  code: '',
  label: ''
}

const { model, mode, makeRequest } = useForm(props, {
  initial,
  updateTransformer({ id, $schema, meta, ...rest }): FixativeUpdate {
    return rest
  }
})

const { field, errorHandler } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? $FixativeInput : $FixativeUpdate)
)

const [loading, toggleLoading] = useToggle(false)

async function submit() {
  toggleLoading(true)
  return await makeRequest({
    create: SamplingService.createFixative,
    edit: ({ code }, body) => SamplingService.updateFixative({ path: { code }, body })
  })
    .then(errorHandler)
    .then((item) => emit('success', item))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped lang="scss"></style>
