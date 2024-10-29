<template>
  <FormDialog title="Register program" v-model="open" v-bind="$attrs" @submit="submit">
    <v-form>
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
          <v-textarea
            label="Description"
            v-model="model.description"
            v-bind="field('description')"
          />
        </v-col>
      </v-row>
    </v-form>
  </FormDialog>
</template>

<script setup lang="ts">
import { $ProgramInput, EventsService, PeopleService, Program, ProgramInput } from '@/api'
import PersonPicker from '../people/PersonPicker.vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { useSchema } from '../toolkit/forms/schema'
import { ref } from 'vue'
import InstitutionPicker from '../people/InstitutionPicker.vue'
import { FormEmits, FormProps, useForm } from '../toolkit/forms/form'
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
const { model, mode } = useForm(props, { initial })

async function submit() {
  const body = model.value
  const request = props.edit
    ? EventsService.updateProgram({ path: { code: props.edit.code }, body })
    : EventsService.createProgram({ body })
  return await request.then(errorHandler).then((program) => emit('success', program))
}
</script>

<style scoped></style>
