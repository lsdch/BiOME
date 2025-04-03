<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} program`"
    @submit="emit('submit', model)"
  >
    <v-container>
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
        <v-col class="d-flex">
          <v-number-input
            label="Start year"
            v-model="model.start_year"
            v-bind="schema('start_year')"
            rounded="e-0"
          />
          <v-number-input
            label="End year"
            v-model="model.end_year"
            v-bind="schema('end_year')"
            rounded="s-0"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <OrganisationPicker
            label="Funding agencies"
            v-model="model.funding_agencies"
            chips
            closable-chips
            item-value="code"
            multiple
            clearable
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <PersonPicker
            label="Managers"
            v-model="model.managers"
            v-bind="schema('managers')"
            item-value="alias"
            multiple
          />
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
import { FormProps } from '@/functions/mutations'
import { ProgramModel } from '@/models'
import FormDialog, { FormDialogProps } from '../toolkit/forms/FormDialog.vue'
import { reactiveComputed } from '@vueuse/core'
import { useSchema } from '@/composables/schema'
import { $ProgramInput, $ProgramUpdate } from '@/api'
import PersonPicker from '../people/PersonPicker.vue'
import OrganisationPicker from '../people/OrganisationPicker.vue'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<ProgramModel.ProgramModel>({
  default: ProgramModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: ProgramModel.ProgramModel | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $ProgramInput : $ProgramUpdate))
</script>

<style scoped lang="scss"></style>
