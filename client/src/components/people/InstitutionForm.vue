<template>
  <v-form ref="innerForm" @submit.prevent="submit">
    <v-container fluid>
      <v-row>
        <v-col>
          <v-text-field
            name="institution"
            label="Institution name"
            id="institution-input"
            hint="The name of the structure for which this instance is deployed, e.g. the name of your lab."
            persistent-hint
            v-model="inst.name"
            required
            :error-messages="errorMsgs.name"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" sm="6">
          <v-text-field
            name="institution_shortname"
            label="Code or abbreviated name"
            id="institution-shortname"
            hint="A short label or code that identifies your lab."
            persistent-hint
            v-model="inst.code"
            :error-messages="errorMsgs.code"
          />
        </v-col>
        <v-col cols="12" sm="6">
          <v-select
            :items="institutionKindOptions"
            v-model="inst.kind"
            label="Kind"
            variant="outlined"
            :itemProps="(item) => ({ title: enumAsString(item) })"
          >
            <template v-slot:prepend-inner>
              <v-icon v-bind="kindIcon(inst.kind)" />
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
            v-model="inst.description"
            :error-messages="errorMsgs.description"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-spacer />
        <v-btn :loading="loading" color="primary" variant="plain" type="submit" text="Submit" />
      </v-row>
    </v-container>
  </v-form>
</template>

<script lang="ts">
const DEFAULT: InstitutionInput = { code: '', name: '', kind: 'Lab' }
</script>

<script setup lang="ts">
import { Institution, InstitutionInput, PeopleService } from '@/api'
import { Ref, ref, watchEffect } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'
import { VForm } from 'vuetify/components'
import { enumAsString } from '../toolkit/enums'
import { Emits, Props, useForm } from '../toolkit/form'
import { institutionKindOptions, kindIcon } from './institutionKind'

const innerForm = ref<ComponentExposed<typeof VForm> | null>(null)

const props = defineProps<Props<Institution>>()
const inst: Ref<InstitutionInput> = ref(DEFAULT)
const emit = defineEmits<Emits<Institution>>()

watchEffect(() => {
  if (props.edit) Object.assign(inst.value, props.edit)
  else Object.assign(inst.value, DEFAULT)
})

function request() {
  if (props.edit) {
    return PeopleService.updateInstitution(props.edit.code, inst.value)
  } else {
    return PeopleService.createInstitution(inst.value)
  }
}

const { errorMsgs, submit, loading } = useForm<InstitutionInput, Institution>(props, emit, request)

function reset() {
  innerForm.value?.reset()
}

defineExpose({
  reset,
  submit
})
</script>

<style scoped></style>
