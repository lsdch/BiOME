<template>
  <FormDialog v-model="dialog" :title="title" :loading="loading" @submit="submit">
    <v-container fluid>
      <v-row>
        <v-col>
          <v-text-field
            id="institution-input"
            v-model="model.name"
            name="institution"
            label="Institution name"
            persistent-hint
            required
            v-bind="field('name')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" sm="6">
          <v-text-field
            id="institution-shortname"
            v-model="model.code"
            name="institution_shortname"
            label="Code or abbreviated name"
            persistent-hint
            v-bind="field('code')"
          />
        </v-col>
        <v-col cols="12" sm="6">
          <v-select
            v-model="model.kind"
            :items="institutionKindOptions"
            v-bind="field('kind')"
            label="Kind"
            variant="outlined"
            :item-props="(item) => ({ title: enumAsString(item) })"
          >
            <template #prepend-inner>
              <v-icon v-bind="kindIcon(model.kind ?? undefined)" />
            </template>
            <template #item="{ item, props }">
              <v-list-item v-bind="props">
                <template #prepend>
                  <v-icon v-bind="kindIcon(item.value)" />
                </template>
              </v-list-item>
            </template>
          </v-select>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea
            v-model="model.description"
            variant="outlined"
            label="Description (optional)"
            v-bind="field('description')"
          />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import {
  $InstitutionInput,
  Institution,
  InstitutionInput,
  InstitutionUpdate,
  PeopleService
} from '@/api'
import { useToggle } from '@vueuse/core'
import { computed } from 'vue'
import { enumAsString } from '../toolkit/enums'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { useForm, useSchema, type FormEmits, type FormProps } from '../toolkit/forms/form'
import { institutionKindOptions, kindIcon } from './institutionKind'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<Institution>>()
const emit = defineEmits<FormEmits<Institution>>()

const title = computed(() => (props.edit ? `Edit ${props.edit.code}` : 'Create institution'))

const [loading, toggleLoading] = useToggle(false)

const initial: InstitutionInput = { code: '', name: '', kind: 'Lab' }
const { model, makeRequest } = useForm<Institution, InstitutionInput, InstitutionUpdate>(props, {
  initial,
  updateTransformer({ code, name, kind, description }: Institution): InstitutionUpdate {
    return { code, name, kind, description }
  }
})

const { errorHandler, field } = useSchema($InstitutionInput)

async function submit() {
  toggleLoading(true)
  await makeRequest({
    create: PeopleService.createInstitution,
    edit: ({ code }, model) => PeopleService.updateInstitution({ path: { code }, body: model })
  })
    .then(errorHandler)
    .then((inst) => emit('success', inst))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped></style>
