<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} person`"
    @submit="emit('submit', model)"
  >
    <v-container fluid>
      <v-row>
        <v-col cols="12" sm="6">
          <v-text-field
            name="first_name"
            label="First name(s)"
            v-model.trim="model.first_name"
            v-bind="schema('first_name')"
          />
        </v-col>
        <v-col cols="12" sm="6">
          <v-text-field
            name="last_name"
            label="Last name"
            v-model.trim="model.last_name"
            v-bind="schema('last_name')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field
            label="Contact (optional)"
            v-model.trim="model.contact"
            prepend-inner-icon="mdi-at"
            v-bind="schema('contact')"
            hint="An e-mail address to contact this person"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <OrganisationPicker
            label="Organisations (optional)"
            v-model="model.organisations"
            item-color="primary"
            chips
            closable-chips
            multiple
            item-value="code"
            v-bind="schema('organisations')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea
            v-model.trim="model.comment"
            variant="outlined"
            label="Comments (optional)"
            v-bind="schema('comment')"
          />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { FormProps } from '@/functions/mutations'
import { PersonModel } from '@/models'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { reactiveComputed } from '@vueuse/core'
import { useSchema } from '@/composables/schema'
import { $PersonInput, $PersonUpdate } from '@/api'
import OrganisationPicker from '@/components/people/OrganisationPicker.vue'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<PersonModel.PersonFormModel>({
  default: PersonModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: PersonModel.PersonFormModel | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $PersonInput : $PersonUpdate))
</script>

<style scoped lang="scss"></style>
