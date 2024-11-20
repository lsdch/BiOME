<template>
  <FormDialog title="Register program" v-model="open" v-bind="$attrs" @submit="submit" :loading>
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
      <v-col class="d-flex">
        <v-number-input
          label="Start year"
          v-model="model.start_year"
          v-bind="field('start_year')"
          rounded="e-0"
        />
        <v-number-input
          label="End year"
          v-model="model.end_year"
          v-bind="field('end_year')"
          rounded="s-0"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <InstitutionPicker
          label="Funding agencies"
          v-model="model.funding_agencies"
          chips
          closable-chips
          item-value="code"
          multiple
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <PersonPicker
          label="Managers"
          v-model="model.managers"
          v-bind="field('managers')"
          multiple
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-textarea label="Description" v-model="model.description" v-bind="field('description')" />
      </v-col>
    </v-row>
  </FormDialog>
</template>

<script setup lang="ts">
import { $ProgramInput, EventsService, Program, ProgramInput, ProgramUpdate } from '@/api'
import { useToggle } from '@vueuse/core'
import InstitutionPicker from '../people/InstitutionPicker.vue'
import PersonPicker from '../people/PersonPicker.vue'
import { FormEmits, FormProps, useForm } from '../toolkit/forms/form'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { useSchema } from '../toolkit/forms/schema'
const open = defineModel<boolean>('open')

const props = defineProps<FormProps<Program>>()
const emit = defineEmits<FormEmits<Program>>()

const initial: ProgramInput = {
  label: '',
  code: '',
  managers: [],
  funding_agencies: []
}

const { field, errorHandler } = useSchema($ProgramInput)
const { model, makeRequest } = useForm(props, {
  initial,
  updateTransformer({ funding_agencies, managers, meta, $schema, ...rest }): ProgramUpdate {
    return {
      ...rest,
      funding_agencies: funding_agencies.map(({ code }) => code),
      managers: managers.map(({ alias }) => alias)
    }
  }
})

const [loading, toggleLoading] = useToggle(false)
async function submit() {
  toggleLoading(true)
  return await makeRequest({
    create: EventsService.createProgram,
    edit: ({ code }, model) => EventsService.updateProgram({ path: { code }, body: model })
  })
    .then(errorHandler)
    .then((program) => emit('success', program))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped></style>
