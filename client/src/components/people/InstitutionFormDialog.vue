<template>
  <FormDialog v-model="dialog" :title="title" :loading="loading" @submit="submit">
    <v-form @submit.prevent="submit">
      <v-container fluid>
        <v-row>
          <v-col>
            <v-text-field
              name="institution"
              label="Institution name"
              id="institution-input"
              persistent-hint
              required
              v-model="model.name"
              v-bind="field('name')"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" sm="6">
            <v-text-field
              name="institution_shortname"
              label="Code or abbreviated name"
              id="institution-shortname"
              persistent-hint
              v-model="model.code"
              v-bind="field('code')"
            />
          </v-col>
          <v-col cols="12" sm="6">
            <v-select
              :items="institutionKindOptions"
              v-model="model.kind"
              v-bind="field('kind')"
              label="Kind"
              variant="outlined"
              :itemProps="(item) => ({ title: enumAsString(item) })"
            >
              <template v-slot:prepend-inner>
                <v-icon v-bind="kindIcon(model.kind)" />
              </template>
              <template v-slot:item="{ item, props }">
                <v-list-item v-bind="props">
                  <template v-slot:prepend>
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
              variant="outlined"
              label="Description (optional)"
              v-model="model.description"
              v-bind="field('description')"
            />
          </v-col>
        </v-row>
      </v-container>
    </v-form>
  </FormDialog>
</template>

<script lang="ts">
const DEFAULT: InstitutionInput = { code: '', name: '', kind: 'Lab' }
</script>

<script setup lang="ts">
import {
  $Institution,
  $InstitutionInput,
  Institution,
  InstitutionInput,
  PeopleService
} from '@/api'
import { computed } from 'vue'
import { VForm } from 'vuetify/components'
import { enumAsString } from '../toolkit/enums'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { useForm, useSchema, type FormEmits, type FormProps } from '../toolkit/forms/form'
import { institutionKindOptions, kindIcon } from './institutionKind'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<Institution>>()
const emit = defineEmits<FormEmits<Institution>>()

const title = computed(() => (props.edit ? `Edit ${props.edit.code}` : 'Create institution'))

async function submit() {
  const req = props.edit
    ? PeopleService.updateInstitution({ path: { code: props.edit.code }, body: model.value })
    : PeopleService.createInstitution({ body: model.value })
  await req.then(errorHandler).then((inst) => emit('success', inst))
}

const { loading, model } = useForm(props, { initial: DEFAULT })
const { errorHandler, field } = useSchema($InstitutionInput)

defineExpose({
  submit
})
</script>

<style scoped></style>
