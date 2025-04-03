<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} organisation`"
    @submit="emit('submit', model)"
  >
    <v-container fluid>
      <v-row>
        <v-col>
          <v-text-field
            id="organisation-input"
            v-model="model.name"
            name="organisation"
            label="Organisation name"
            persistent-hint
            required
            v-bind="schema('name')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" sm="6">
          <v-text-field
            id="organisation-shortname"
            v-model="model.code"
            name="organisation_shortname"
            label="Code or abbreviated name"
            persistent-hint
            v-bind="schema('code')"
          />
        </v-col>
        <v-col cols="12" sm="6">
          <OrgKindPicker
            v-model="model.kind"
            v-bind="schema('kind')"
            label="Kind"
            variant="outlined"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea
            v-model="model.description"
            variant="outlined"
            label="Description (optional)"
            v-bind="schema('description')"
          />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import {
  $OrganisationInput,
  $OrganisationUpdate,
  OrganisationInput,
  OrganisationUpdate
} from '@/api'
import OrgKindPicker from '@/components/people/OrgKindPicker.vue'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { useSchema } from '@/composables/schema'
import { FormProps } from '@/functions/mutations'
import { OrganisationModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<OrganisationInput | OrganisationUpdate>({
  default: OrganisationModel.initialModel
})
const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: OrganisationInput | OrganisationUpdate | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $OrganisationInput : $OrganisationUpdate))
</script>

<style scoped></style>
